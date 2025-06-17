package tags

import (
	"smart-scene-app-api/server"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(ctx server.ServerContext, router *gin.RouterGroup) {
	tagHandler := NewTagHandler(ctx)

	// Tag routes
	tagRoutes := router.Group("/tags")
	{
		// GET /api/tags/hierarchy - Lấy hierarchical tags cho sidebar filters
		tagRoutes.GET("/hierarchy", tagHandler.GetTagsHierarchy)

		// GET /api/tags - Lấy flat list tags với pagination
		tagRoutes.GET("", tagHandler.GetTags)

		// GET /api/tags/position/:position_code - Tags theo position
		tagRoutes.GET("/position/:position_code", tagHandler.GetTagsByPosition)

		// GET /api/tags/category/:category_code - Tags theo category
		tagRoutes.GET("/category/:category_code", tagHandler.GetTagsByCategory)

		// GET /api/tags/stats - Thống kê usage của tags
		tagRoutes.GET("/stats", tagHandler.GetTagUsageStats)
	}
}
