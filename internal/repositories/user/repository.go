package user

import (
	"smart-scene-app-api/internal/models/user"
	"smart-scene-app-api/internal/repositories"

	"gorm.io/gorm"
)

// Repository defines the user repository interface
type Repository interface {
	repositories.BaseRepository[user.User]
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) Repository {
	return repositories.NewBaseRepository[user.User](db)
}
