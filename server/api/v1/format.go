package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/txfs19260817/scopelens/server/models"
	"github.com/txfs19260817/scopelens/server/utils/response"
)

func GetFormats(c *gin.Context) {
	response.OkWithData(models.GetFormats(), c)
}
