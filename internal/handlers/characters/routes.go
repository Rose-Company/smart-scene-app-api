package character

import (
	"smart-scene-app-api/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		videos := v1.Group("/videos")
		{
			videos.GET("/:id/characters", middleware.UserAuthentication(), h.GetCharactersByVideoID)
		}
	}
}
