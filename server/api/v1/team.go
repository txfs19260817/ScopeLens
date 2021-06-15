package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/config"
	"github.com/txfs19260817/scopelens/server/models"
	"github.com/txfs19260817/scopelens/server/utils/response"
	"github.com/txfs19260817/scopelens/server/utils/validator"
	"github.com/unknwon/com"
	"go.uber.org/zap"
)

func InsertTeam(c *gin.Context) {
	var team models.Team

	// Validate Team form
	if err, valid := validator.TeamValidator(&team, c.Request); !valid {
		zap.L().Error("validating team error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// Insert a team
	if _, err := models.Db.InsertTeam(team); err != nil {
		zap.L().Error("inserting team error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Inserted a team successfully! ", c)
	}
}

func GetTeams(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		zap.L().Warn("parsing parameter `page` error, reset it to 0", zap.Error(err))
		page = 0
	}

	// retrieve data
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, models.Search{OrderBy: "time"}, false); err != nil {
		zap.L().Error("get teams error", zap.Error(err), zap.String("order", "time"))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetTeamsOrderByLikes(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		zap.L().Warn("parsing parameter `page` error, reset it to 0", zap.Error(err))
		page = 0
	}

	// retrieve data
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, models.Search{OrderBy: "likes"}, false); err != nil {
		zap.L().Error("get teams error", zap.Error(err), zap.String("order", "likes"))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetTeamByID(c *gin.Context) {
	team, err := models.Db.GetTeamByID(c.Param("id"))
	if err != nil {
		zap.L().Error("get team by ID error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(team, c)
	}
}

func GetTeamsBySearchCriteria(c *gin.Context) {
	page, err := com.StrTo(c.Query("page")).Int()
	if err != nil {
		zap.L().Warn("parsing parameter `page` error, reset it to 0", zap.Error(err))
		page = 0
	}

	// Validate JSON form
	var s models.Search
	if err := c.ShouldBindJSON(&s); err != nil {
		zap.L().Error("decoding search criteria error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// retrieve data
	data := make(map[string]interface{})
	if data["teams"], data["total"], err = models.Db.GetTeams(page, config.App.PageSize, s, true); err != nil {
		zap.L().Error("search teams error", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}

func GetPokemonUsageByFormat(c *gin.Context) {
	format := c.Param("format")
	usages, err := models.Db.GetPokemonUsageByFormat(format)
	if err != nil {
		zap.L().Error("get usage error", zap.Error(err))
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
		zap.L().Error("get liked teams by username error", zap.Error(err), zap.String("by", "liked"))
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
		zap.L().Error("get uploaded teams by username error", zap.Error(err), zap.String("by", "uploaded"))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(data, c)
	}
}
