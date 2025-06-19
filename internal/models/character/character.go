package character

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

type Character struct {
	models.Base
	Name        string      `json:"name" gorm:"type:text;not null;index"`
	Description string      `json:"description" gorm:"type:text"`
	Avatar      string      `json:"avatar" gorm:"type:text"`
	Metadata    common.JSON `json:"metadata" gorm:"type:jsonb"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedBy   uuid.UUID   `json:"created_by" gorm:"type:uuid;not null"`
	UpdatedBy   uuid.UUID   `json:"updated_by" gorm:"type:uuid;not null"`
}

type CharacterFilterAndPagination struct {
	Name        string             `json:"name" form:"name"`
	IsActive    *bool              `json:"is_active" form:"is_active"`
	CreatedBy   uuid.UUID          `json:"created_by" form:"created_by"`
	QueryParams models.QueryParams `json:"query_params" form:"query_params"`
}
