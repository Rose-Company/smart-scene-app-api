package character

import (
	"smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"

	"gorm.io/gorm"
)

type AppearanceRepository interface {
	repositories.BaseRepository[character.CharacterAppearance]
}

func NewAppearanceRepository(db *gorm.DB) AppearanceRepository {
	return repositories.NewBaseRepository[character.CharacterAppearance](db)
}
