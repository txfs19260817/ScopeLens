package routers

import (
	"github.com/gin-gonic/gin"
	v1 "scopelens-server/api/v1"
)

func InitTeamRouter(Router *gin.RouterGroup) {
	TeamRouter := Router.Group("team")//.Use(middleware.JWTAuth())
	{
		TeamRouter.POST("/post", v1.InsertTeam)
	}
}