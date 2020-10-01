package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"scopelens-server/middleware"
	"scopelens-server/models"
	"scopelens-server/utils/logger"
	util "scopelens-server/utils/recaptcha"
	"scopelens-server/utils/response"
)

func Register(c *gin.Context) {
	if err := util.ReCaptcha(c.Query("recaptcha")); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	var user models.User
	// Validate JSON form
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check if the username already exists.
	legalName, err := models.Db.CheckUsernameAvailability(user.UserName)
	if err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Existed
	if !legalName {
		err := fmt.Errorf("Username already exists. ")
		logger.SugaredLogger.Error(err)
		response.FailWithMessage("Username already exists. ", c)
		return
	}
	// Insert new user
	if _, err := models.Db.Register(user); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Registered successfully! ", c)
	}
}

func Login(c *gin.Context) {
	var loginReq models.Login
	// Validate JSON form
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check username and password
	_, err := models.Db.LoginValidate(loginReq)
	if err != nil {
		logger.SugaredLogger.Error(err)
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
