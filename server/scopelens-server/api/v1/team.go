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
	data := make(map[string]interface{})

	// retrieve data
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, "time", "", []string{}); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetTeamsOrderbyLikes(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		page = 0
	}
	// Get the total number of teams
	data := make(map[string]interface{})

	// retrieve data
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, "likes", "", []string{}); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetTeamByID(c *gin.Context) {
	id := c.Param("id")
	team, err := models.Db.GetTeamByID(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(team, c)
	}
}

func GetTeamsBySearchCriteria(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		page = 0
	}

	// define search criteria struct
	type Search struct {
		Format  string   `bson:"format" json:"format"`
		Pokemon []string `bson:"pokemon" json:"pokemon"`
	}

	// Validate JSON form
	var s Search
	if err := c.ShouldBindJSON(&s); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Get the total number of teams
	data := make(map[string]interface{})

	// retrieve data
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, "time", s.Format, s.Pokemon); err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}
