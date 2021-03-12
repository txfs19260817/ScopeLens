package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/txfs19260817/scopelens/server/api/v1"
	"github.com/txfs19260817/scopelens/server/middleware"
)

func InitAuthRouter(Router *gin.RouterGroup) {
	AuthRouter := Router.Group("auth")
	{
		AuthRouter.POST("/register", v1.Register)
		AuthRouter.POST("/login", v1.Login)
	}
}

func InitAuthAUTHRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("auth").Use(middleware.JWTAuth())
	{
		UserRouter.GET("/checktoken", v1.CheckToken)
	}
}
