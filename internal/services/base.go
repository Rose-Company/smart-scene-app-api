package services

import (
	"smart-scene-app-api/internal/services/auth"
	"smart-scene-app-api/internal/services/video"
	l "smart-scene-app-api/pkg/logger"
	"smart-scene-app-api/server"

	"go.uber.org/zap"
)

type Service struct {

	Auth auth.Service
	Video video.Service
	logger *zap.Logger

}

func NewService(sc server.ServerContext) *Service {

	authService := auth.NewAuthService(sc)
	videoService := video.NewVideoService(sc)
	return &Service{
		logger: l.New(),
		Auth:   authService,
		Video:  videoService,
	}
}

