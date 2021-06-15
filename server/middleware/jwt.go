package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/config"
	"github.com/txfs19260817/scopelens/server/utils/response"
)

var (
	TokenExpired     = fmt.Errorf("Token is expired. ")
	TokenNotValidYet = fmt.Errorf("Token not active yet. ")
	TokenMalformed   = fmt.Errorf("Illegal token. ")
	TokenInvalid     = fmt.Errorf("Could not handle this token: ")
)

// JWTAuth middleware, validate the token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": -1,
				"data": nil,
				"msg":  "Token not found. ",
			})
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": nil,
				"msg":  "Bad authorization format, correct format：`Authorization: Bearer <token>`. ",
			})
			c.Abort()
			return
		}
		token := parts[1]
		j := NewJWT()
		// parse token
		claims, err := j.ParserToken(token)

		if err != nil {
			// token过期
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": -1,
					"data": nil,
					"msg":  err.Error(),
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": nil,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		// 解析到具体的claims相关信息
		c.Set("claims", claims)
		c.Next()
	}
}

// JWT基本数据结构
// 签名的signkey
type JWT struct {
	SigningKey []byte
}

// 定义载荷
type CustomClaims struct {
	UserName string `json:"userName"`
	// Email    string `json:"email"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

// 初始化JWT实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// Get signKey
func GetSignKey() string {
	return config.Jwt.SigningKey
}

// 创建Token(基于用户的基本信息claims)
// 使用HS256算法进行token生成
// 使用用户基本信息claims以及签名key(signkey)生成token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	// return a token struct pointer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// token parsing
// Couldn't handle this token:
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	// 将token中的claims信息解析出来和用户原始数据进行校验
	// 做以下类型断言，将token.Claims转换成具体用户自定义的Claims结构体
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid

}

// 更新Token
func (j *JWT) UpdateToken(tokenString string) (string, error) {
	// TimeFunc为一个默认值是time.Now的当前时间变量,用来解析token后进行过期时间验证
	// 可以使用其他的时间值来覆盖
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	// 拿到token基础数据
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil

	})

	// 校验token当前还有效
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// 修改Claims的过期时间(int64)
		// https://gowalker.org/github.com/dgrijalva/jwt-go#StandardClaims
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", err
}

// token generator
func GenerateToken(c *gin.Context, username string) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := NewJWT()

	// 构造用户claims信息(payload)
	claims := CustomClaims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),        // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 60*60*24*30), // 签名过期时间30days
			Issuer:    "ZeminJiang",                           // 签名颁发者
		},
	}

	// 根据claims生成token对象
	token, err := j.CreateToken(claims)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		c.Abort()
	} else {
		response.OkWithData(token, c)
	}
}
