package v1

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"scopelens-server/middleware"
	"scopelens-server/models"
	"scopelens-server/utils/response"
	"time"
)

func Register(c *gin.Context) {
	var user models.User
	// Validate JSON form
	if err := c.ShouldBindJSON(&user); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check if the username already exists.
	legalName, err := models.Db.CheckUsernameAvailability(user.UserName)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Existed
	if !legalName {
		response.FailWithMessage("Username already exists. ", c)
		return
	}
	// Insert new user
	if _, err := models.Db.Register(user); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Registered successfully! ", c)
	}
}

func Login(c *gin.Context) {
	var loginReq models.Login
	// Validate JSON form
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check username and password
	_, err := models.Db.LoginValidate(loginReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		// respond a token
		generateToken(c, loginReq.UserName)
	}
}

// token generator
func generateToken(c *gin.Context, username string) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := middleware.NewJWT()

	// 构造用户claims信息(payload)
	claims := middleware.CustomClaims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),          // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + + 60*60*24*30), // 签名过期时间30days
			Issuer:    "ZeminJiang",                             // 签名颁发者
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
