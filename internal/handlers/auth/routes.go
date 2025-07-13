package auth

import (
	"smart-scene-app-api/internal/services"
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

	public := router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/register", h.Register)
		}
	}
}
