package tags

import (
	"smart-scene-app-api/middleware"
	"smart-scene-app-api/server"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(ctx server.ServerContext, router *gin.RouterGroup) {
	tagHandler := NewTagHandler(ctx)

	tagRoutes := router.Group("/tags")
	{
		tagRoutes.GET("/position/:position_code", middleware.UserAuthentication(), tagHandler.GetTagsByPosition)
	}
}
