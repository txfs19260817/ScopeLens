package routers

import (
	"github.com/gin-gonic/gin"
	"scopelens-server/middleware"
)

func InitRouters() *gin.Engine {
	r := gin.Default()

	// static resources
	r.Static("/assets", "./assets")

	// middleware
	r.Use(middleware.Cors())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiGroups := r.Group("api")
	InitUserRouter(apiGroups)
	InitAuthRouter(apiGroups)
	InitAuthAUTHRouter(apiGroups)
	InitFormatRouter(apiGroups)
	InitTeamRouter(apiGroups)
	InitTeamAUTHRouter(apiGroups)

	return r
}
