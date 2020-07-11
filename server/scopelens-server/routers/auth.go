package routers

import (
	"github.com/gin-gonic/gin"
	v1 "scopelens-server/api/v1"
)

func InitAuthRouter(Router *gin.RouterGroup) {
	AuthRouter := Router.Group("auth")
	{
		AuthRouter.POST("/register", v1.Register)
		AuthRouter.POST("/login", v1.Login)
	}
}