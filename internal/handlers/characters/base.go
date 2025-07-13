package character

import (
	"smart-scene-app-api/internal/services"
	"smart-scene-app-api/server"

	"go.uber.org/zap"
)

type Handler struct {
	sc      server.ServerContext
	service *services.Services
	logger  *zap.Logger
}

func NewHandler(sc server.ServerContext) *Handler {
	return &Handler{
		sc:      sc,
		service: services.NewServices(sc),
		logger:  zap.NewExample(),
	}
}
