package services

import (
	l "smart-scene-app-api/pkg/logger"
	"smart-scene-app-api/server"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func NewService(sc server.ServerContext) *Service {

	return &Service{
		logger: l.New(),
	}
}
