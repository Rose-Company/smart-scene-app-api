package handlers

import (
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

func (h *Handler) RegisterRouter(c *gin.Engine) {
	// authConfig := h.sc.GetAuthConfig()
	// authenticator := middleware.NewAuthenticator(authConfig)

}
