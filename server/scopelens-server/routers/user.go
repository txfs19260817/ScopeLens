package routers

import (
	"github.com/gin-gonic/gin"
	v1 "scopelens-server/api/v1"
	"scopelens-server/middleware"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(middleware.JWTAuth())
	{
		UserRouter.GET("/username/:username", v1.GetUserByName)
	}
}