package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/txfs19260817/scopelens/server/api/v1"
	"github.com/txfs19260817/scopelens/server/middleware"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").Use(middleware.JWTAuth())
	{
		UserRouter.GET("/username/:username", v1.GetUserByName)
		UserRouter.POST("/like", v1.InsertLikeByUsername)
	}
}
