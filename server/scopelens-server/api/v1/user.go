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
		response.FailWithMessage("GetUserByName Fail: " + err.Error(), c)
	} else {
		response.OkWithData(user, c)
	}
}
