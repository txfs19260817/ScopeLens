package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/middleware"
	"github.com/txfs19260817/scopelens/server/models"
	util "github.com/txfs19260817/scopelens/server/utils/recaptcha"
	"github.com/txfs19260817/scopelens/server/utils/response"
	"go.uber.org/zap"
)

func Register(c *gin.Context) {
	if err := util.ReCaptcha(c.Query("recaptcha")); err != nil {
		zap.L().Error("recaptcha error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	var user models.User
	// Validate JSON form
	if err := c.ShouldBindJSON(&user); err != nil {
		zap.L().Error("decoding register data error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check if the username already exists.
	legalName, err := models.Db.CheckUsernameAvailability(user.UserName)
	if err != nil {
		zap.L().Error("validating username error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Existed
	if !legalName {
		err := fmt.Errorf("the username already exists. ")
		zap.L().Warn("illegal username", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Insert new user
	if _, err := models.Db.Register(user); err != nil {
		zap.L().Error("insert username error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Registered successfully! ", c)
	}
}

func Login(c *gin.Context) {
	var loginReq models.Login
	// Validate JSON form
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		zap.L().Error("decoding login data error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check username and password
	_, err := models.Db.LoginValidate(loginReq)
	if err != nil {
		zap.L().Error("check username and password error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		// respond a token
		middleware.GenerateToken(c, loginReq.UserName)
	}
}

func CheckToken(c *gin.Context) {
	// api combined with JWTAuth middleware to check if token is still valid
	response.Ok(c)
}
