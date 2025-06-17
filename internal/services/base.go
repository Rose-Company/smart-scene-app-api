package services

import (
	"smart-scene-app-api/internal/services/auth"
	"smart-scene-app-api/internal/services/character"
	"smart-scene-app-api/internal/services/tag"
	"smart-scene-app-api/internal/services/video"
	l "smart-scene-app-api/pkg/logger"
	"smart-scene-app-api/server"

	"go.uber.org/zap"
)

type Service struct {
	Auth      auth.Service
	Video     video.Service
	Character character.Service
	Tag       tag.Service
	logger    *zap.Logger
}

// Services alias for consistency
type Services = Service

func NewService(sc server.ServerContext) *Service {
	return NewServices(sc)
}

func NewServices(sc server.ServerContext) *Services {

	authService := auth.NewAuthService(sc)
	videoService := video.NewVideoService(sc)
	characterService := character.NewCharacterService(sc)
	tagService := tag.NewTagService(sc)

	return &Services{
		logger:    l.New(),
		Auth:      authService,
		Video:     videoService,
		Character: characterService,
		Tag:       tagService,
	}
}
