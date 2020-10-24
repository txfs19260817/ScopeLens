package v1

import (
	"github.com/gin-gonic/gin"
	"scopelens-server/models"
	"scopelens-server/utils/response"
)

func GetFormats(c *gin.Context) {
	response.OkWithData(models.GetFormats(), c)
}
