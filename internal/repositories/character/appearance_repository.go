package character

import (
	"smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"

	"gorm.io/gorm"
)

type AppearanceRepository interface {
	repositories.BaseRepository[character.CharacterAppearance]
}

type appearanceRepository struct {
	repositories.BaseRepository[character.CharacterAppearance]
	db *gorm.DB
}

func NewAppearanceRepository(db *gorm.DB) AppearanceRepository {
	baseRepo := repositories.NewBaseRepository[character.CharacterAppearance](db)
	return &appearanceRepository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
