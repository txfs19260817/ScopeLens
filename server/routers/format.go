package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/txfs19260817/scopelens/server/api/v1"
)

func InitFormatRouter(Router *gin.RouterGroup) {
	FormatRouter := Router.Group("format")
	{
		FormatRouter.GET("/", v1.GetFormats)
	}
}
