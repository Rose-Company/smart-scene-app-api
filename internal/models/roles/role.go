package roles

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

type Role struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name     string    `gorm:"type:text;not null" json:"name"`
	PublicID string    `gorm:"type:text" json:"public_id"`
	models.Base
}

func (Role) TableName() string {
	return common.POSTGRES_TABLE_NAME_ROLES
}
