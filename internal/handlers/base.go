package handlers

import (
	authHandler "smart-scene-app-api/internal/handlers/auth"
	characterHandler "smart-scene-app-api/internal/handlers/characters"
	tagHandler "smart-scene-app-api/internal/handlers/tags"
	videoHandler "smart-scene-app-api/internal/handlers/videos"
	services "smart-scene-app-api/internal/services"
	l "smart-scene-app-api/pkg/logger"
	"smart-scene-app-api/server"

	"github.com/gin-gonic/gin"
)

var ll = l.New()

type Handler struct {
	sc      server.ServerContext
	service *services.Service
}

func NewHandler(sc server.ServerContext) *Handler {
	return &Handler{
		sc:      sc,
		service: services.NewService(sc),
	}
}

func (h *Handler) RegisterRouter(router *gin.Engine) {
	auth := authHandler.NewHandler(h.sc)
	auth.RegisterRoutes(router)

	video := videoHandler.NewHandler(h.sc)
	video.RegisterRoutes(router)

	character := characterHandler.NewHandler(h.sc)
	character.RegisterRoutes(router)

	tagRoutes := router.Group("/api/v1")
	tagHandler.RegisterTagRoutes(h.sc, tagRoutes)
}
