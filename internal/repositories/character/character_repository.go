package character

import (
	"smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"

	"gorm.io/gorm"
)

type Repository interface {
	repositories.BaseRepository[character.Character]
}

func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[character.Character](db)
}
