package character

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"
	"time"

	"github.com/google/uuid"
)

type CharacterAppearance struct {
	ID          int         `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time   `json:"created_at" gorm:"type:timestamp;not null;default:now()"`
	CreatedBy   uuid.UUID   `json:"created_by" gorm:"type:uuid;not null;index"`
	VideoID     uuid.UUID   `json:"video_id" gorm:"type:uuid;not null;index"`
	CharacterID uuid.UUID   `json:"character_id" gorm:"type:uuid;not null;index"`
	StartFrame  int         `json:"start_frame" gorm:"not null;index"`
	EndFrame    int         `json:"end_frame" gorm:"not null;index"`
	StartTime   float64     `json:"start_time" gorm:"type:decimal(10,3);not null"`
	EndTime     float64     `json:"end_time" gorm:"type:decimal(10,3);not null"`
	Duration    float64     `json:"duration" gorm:"type:decimal(10,3);default:0"`
	Confidence  float64     `json:"confidence" gorm:"type:decimal(5,4);default:0"`
	Metadata    common.JSON `json:"metadata" gorm:"type:jsonb"`

	Video     interface{} `json:"video,omitempty" gorm:"foreignKey:VideoID"`
	Character *Character  `json:"character,omitempty" gorm:"foreignKey:CharacterID;references:ID"`
}

func (c *CharacterAppearance) TableName() string {
	return common.POSTGRES_TABLE_NAME_CHARACTER_APPEARANCES
}

type VideoCharacterFilterAndPagination struct {
	models.BaseRequestParamsUri
}

type VideoCharacterSummary struct {
	VideoID         uuid.UUID `json:"video_id"`
	CharacterID     uuid.UUID `json:"character_id"`
	CharacterName   string    `json:"character_name"`
	CharacterAvatar string    `json:"character_avatar"`
}

type VideoCharacterListResponse struct {
	models.BaseListResponse
	Items []VideoCharacterSummary `json:"items"`
}
