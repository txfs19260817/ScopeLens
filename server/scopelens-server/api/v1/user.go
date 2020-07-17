package v1

import (
	"github.com/gin-gonic/gin"
	"scopelens-server/models"
	"scopelens-server/utils/response"
)

func GetUserByName(c *gin.Context) {
	username := c.Param("username")
	user, err := models.Db.GetUserByUsername(username)
	if err != nil {
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
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := models.Db.InsertLikeByUsername(l.Username, l.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}
