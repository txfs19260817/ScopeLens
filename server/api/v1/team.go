package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/config"
	"github.com/txfs19260817/scopelens/server/models"
	"github.com/txfs19260817/scopelens/server/utils/logger"
	"github.com/txfs19260817/scopelens/server/utils/response"
	"github.com/txfs19260817/scopelens/server/utils/validator"
	"github.com/unknwon/com"
)

func InsertTeam(c *gin.Context) {
	var team models.Team

	// Validate Team form
	if err, valid := validator.TeamValidator(&team, c.Request); !valid {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Insert a team
	if _, err := models.Db.InsertTeam(team); err != nil {
		logger.SugaredLogger.Error(err)
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
	s := models.Search{OrderBy: "time"}

	// retrieve data
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, s); err != nil {
		logger.SugaredLogger.Error(err)
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
	s := models.Search{OrderBy: "likes"}

	// retrieve data
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, s); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetTeamByID(c *gin.Context) {
	id := c.Param("id")
	team, err := models.Db.GetTeamByID(id)
	if err != nil {
		logger.SugaredLogger.Error(err)
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

	// Validate JSON form
	var s models.Search
	if err := c.ShouldBindJSON(&s); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	// retrieve data
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, s); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetPokemonUsageByFormat(c *gin.Context) {
	format := c.Param("format")
	usages, err := models.Db.GetPokemonUsageByFormat(format)
	if err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(usages, c)
	}
}

func GetLikedTeamsByUsername(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		page = 0
	}

	username := c.Param("username")
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetLikedTeamsByUsername(page, config.App.PageSize, username); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetUploadedTeamsByUsername(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		page = 0
	}

	username := c.Param("username")
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetUploadedTeamsByUsername(page, config.App.PageSize, username); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}