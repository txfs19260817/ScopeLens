package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/models"
	"github.com/txfs19260817/scopelens/server/utils/logger"
	"github.com/txfs19260817/scopelens/server/utils/response"
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
