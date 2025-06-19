package character

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

type CharacterAppearance struct {
	models.Base
	VideoID     uuid.UUID   `json:"video_id" gorm:"type:uuid;not null;index"`
	CharacterID uuid.UUID   `json:"character_id" gorm:"type:uuid;not null;index"`
	StartFrame  int         `json:"start_frame" gorm:"not null;index"`
	EndFrame    int         `json:"end_frame" gorm:"not null;index"`
	StartTime   string      `json:"start_time" gorm:"type:text;not null"`
	EndTime     string      `json:"end_time" gorm:"type:text;not null"`
	Duration    float64     `json:"duration" gorm:"type:decimal(10,3);default:0"`
	Metadata    common.JSON `json:"metadata" gorm:"type:jsonb"`

	Video     interface{} `json:"video,omitempty" gorm:"foreignKey:VideoID"`
	Character Character   `json:"character,omitempty" gorm:"foreignKey:CharacterID"`
}

type CharacterAppearanceFilterAndPagination struct {
	VideoID           uuid.UUID          `json:"video_id" form:"video_id"`
	CharacterID       uuid.UUID          `json:"character_id" form:"character_id"`
	CharacterIDs      []uuid.UUID        `json:"character_ids" form:"character_ids"`
	IncludeCharacters []uuid.UUID        `json:"include_characters" form:"include_characters"`
	ExcludeCharacters []uuid.UUID        `json:"exclude_characters" form:"exclude_characters"`
	StartTimeFrom     string             `json:"start_time_from" form:"start_time_from"`
	StartTimeTo       string             `json:"start_time_to" form:"start_time_to"`
	MinDuration       float64            `json:"min_duration" form:"min_duration"`
	MaxDuration       float64            `json:"max_duration" form:"max_duration"`
	MinConfidence     float64            `json:"min_confidence" form:"min_confidence"`
	QueryParams       models.QueryParams `json:"query_params" form:"query_params"`
}

type VideoCharacterSummary struct {
	VideoID         uuid.UUID `json:"video_id"`
	CharacterID     uuid.UUID `json:"character_id"`
	CharacterName   string    `json:"character_name"`
	CharacterAvatar string    `json:"character_avatar"`
	AppearanceCount int       `json:"appearance_count"`
	TotalDuration   float64   `json:"total_duration"`
	FirstAppearance string    `json:"first_appearance"`
	LastAppearance  string    `json:"last_appearance"`
}

type SceneSegment struct {
	VideoID        uuid.UUID   `json:"video_id"`
	CharacterIDs   []uuid.UUID `json:"character_ids"`
	CharacterNames []string    `json:"character_names"`
	StartTime      string      `json:"start_time"`
	EndTime        string      `json:"end_time"`
	Duration       float64     `json:"duration"`
	StartFrame     int         `json:"start_frame"`
	EndFrame       int         `json:"end_frame"`
}
