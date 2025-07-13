package character

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"
	"time"

	"github.com/google/uuid"
)

type Character struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt   time.Time   `json:"created_at" gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"type:timestamp;not null;default:now()"`
	CreatedBy   uuid.UUID   `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedBy   uuid.UUID   `json:"updated_by" gorm:"type:uuid;not null"`
	Name        string      `json:"name" gorm:"type:text;not null;index"`
	Description string      `json:"description" gorm:"type:text"`
	Avatar      string      `json:"avatar" gorm:"type:text"`
	Metadata    common.JSON `json:"metadata" gorm:"type:jsonb"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
}

func (c *Character) TableName() string {
	return common.POSTGRES_TABLE_NAME_CHARACTERS
}

type CharacterFilterAndPagination struct {
	Name        string             `json:"name" form:"name"`
	IsActive    *bool              `json:"is_active" form:"is_active"`
	CreatedBy   uuid.UUID          `json:"created_by" form:"created_by"`
	QueryParams models.QueryParams `json:"query_params" form:"query_params"`
}
