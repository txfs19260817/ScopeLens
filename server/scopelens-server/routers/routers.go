package routers

import (
	"github.com/gin-gonic/gin"
	"scopelens-server/config"
	"scopelens-server/middleware"
)

func InitRouters() *gin.Engine {
	r := gin.Default()

	// static resources
	r.Static("/assets", "./assets")

	// middleware
	r.Use(middleware.Cors())
	if config.Mode == "release" {
		r.Use(middleware.TlsHandler())
	}

	apiGroups := r.Group("api")
	InitUserRouter(apiGroups)
	InitAuthRouter(apiGroups)
	InitAuthAUTHRouter(apiGroups)
	InitFormatRouter(apiGroups)
	InitTeamRouter(apiGroups)
	InitTeamAUTHRouter(apiGroups)

	return r
}
