package v1

import (
	"github.com/gin-gonic/gin"
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
