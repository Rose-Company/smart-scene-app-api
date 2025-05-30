package services

import (
	"smart-scene-app-api/internal/services/auth"
	l "smart-scene-app-api/pkg/logger"
	"smart-scene-app-api/server"

	"go.uber.org/zap"
)

type Service struct {
	Auth auth.Service
	logger *zap.Logger
}

func NewService(sc server.ServerContext) *Service {

	authService := auth.NewAuthService(sc)
	return &Service{
		logger: l.New(),
		Auth:   authService,
	}
}

