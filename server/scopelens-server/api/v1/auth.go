package v1

import (
	"github.com/gin-gonic/gin"
	"scopelens-server/middleware"
	"scopelens-server/models"
	"scopelens-server/utils/response"
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
		middleware.GenerateToken(c, loginReq.UserName)
	}
}

func CheckToken(c *gin.Context)  {
	// api combined with JWTAuth middleware to check if token is still valid
	response.Ok(c)
}
