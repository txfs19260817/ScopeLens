package v1

import (
	"scopelens-server/models"
	"scopelens-server/utils/logger"
	"scopelens-server/utils/response"

	"github.com/gin-gonic/gin"
)

func GetUserByName(c *gin.Context) {
	username := c.Param("username")
	user, err := models.Db.GetUserByUsername(username)
	if err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(user, c)
	}
}

func InsertLikeByUsername(c *gin.Context) {
	// define like struct
	type Like struct {
		Username string `bson:"username" json:"username"`
		ID       string `bson:"_id" json:"id"`
	}

	// Validate JSON form
	var l Like
	if err := c.ShouldBindJSON(&l); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := models.Db.InsertLikeByUsername(l.Username, l.ID); err != nil {
		logger.SugaredLogger.Error(err)
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}
