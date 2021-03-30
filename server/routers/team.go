package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/txfs19260817/scopelens/server/api/v1"
	"github.com/txfs19260817/scopelens/server/middleware"
)

func InitTeamAUTHRouter(Router *gin.RouterGroup) {
	TeamRouter := Router.Group("team").Use(middleware.JWTAuth())
	{
		TeamRouter.POST("/post", v1.InsertTeam)
	}
}

func InitTeamRouter(Router *gin.RouterGroup) {
	TeamRouter := Router.Group("team")
	{
		TeamRouter.GET("/teams", v1.GetTeams)
		TeamRouter.GET("/likes", v1.GetTeamsOrderByLikes)
		TeamRouter.GET("/teams/:id", v1.GetTeamByID)
		TeamRouter.POST("/search", v1.GetTeamsBySearchCriteria)
		TeamRouter.GET("/usage/:format", v1.GetPokemonUsageByFormat)
		TeamRouter.GET("/likes/:username", v1.GetLikedTeamsByUsername)
		TeamRouter.GET("/uploaded/:username", v1.GetUploadedTeamsByUsername)
	}
}
