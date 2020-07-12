package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"scopelens-server/config"
	"scopelens-server/models"
	"scopelens-server/utils/response"
	"scopelens-server/utils/validator"
)

func InsertTeam(c *gin.Context) {
	var team models.Team

	// Validate Team form
	if err, valid := validator.TeamValidator(&team, c.Request); !valid {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Insert a team
	if _, err := models.Db.InsertTeam(team); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Inserted a team successfully! ", c)
	}
}

func GetTeams(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		page = 0
	}
	teams, err := models.Db.GetTeams(page, config.App.PageSize, "time")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(teams, c)
	}
}

func GetTeamsOrderbyLikes(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		page = 0
	}
	teams, err := models.Db.GetTeams(page, config.App.PageSize, "likes")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(teams, c)
	}
}

func GetTeamByID(c *gin.Context)  {
	id := c.Param("id")
	team, err := models.Db.GetTeamByID(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(team, c)
	}
}