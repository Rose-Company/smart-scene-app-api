package video

import (
	"smart-scene-app-api/internal/models/video"
	"smart-scene-app-api/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[video.Video]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[video.Video](db)
}
