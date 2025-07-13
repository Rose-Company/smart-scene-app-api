package role

import (
	roleModels "smart-scene-app-api/internal/models/roles"
	"smart-scene-app-api/internal/repositories"

	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
	repositories.BaseRepository[roleModels.Role]
}

func NewRoleRepository(db *gorm.DB) *RoleRepo {
	baseRepo := repositories.NewBaseRepository[roleModels.Role](db)
	return &RoleRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
