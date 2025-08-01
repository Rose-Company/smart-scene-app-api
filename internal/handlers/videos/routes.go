package video

import (
	"smart-scene-app-api/internal/services"
	"smart-scene-app-api/middleware"
	"smart-scene-app-api/server"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	sc      server.ServerContext
	service *services.Service
	logger  *zap.Logger
}

func NewHandler(sc server.ServerContext) *Handler {
	return &Handler{
		sc:      sc,
		service: services.NewService(sc),
		logger:  zap.NewExample(),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {

	protected := router.Group("/api/v1")
	{
		videos := protected.Group("/videos")
		{
			videos.GET("", middleware.UserAuthentication(), h.GetAllVideos)
			videos.GET("/:id", middleware.UserAuthentication(), h.GetVideoDetail)

			videos.POST("", middleware.UserAuthentication(), h.CreateVideo)
			videos.PUT("/:id", middleware.UserAuthentication(), h.UpdateVideo)
			videos.DELETE("/:id", middleware.UserAuthentication(), h.DeleteVideo)
		}
	}
}
