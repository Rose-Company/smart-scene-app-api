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

type VideoCharacterFilterAndPagination struct {
	models.BaseRequestParamsUri
	CharacterName  string  `json:"character_name" form:"character_name"`
	MinConfidence  float64 `json:"min_confidence" form:"min_confidence"`
	MinAppearances int     `json:"min_appearances" form:"min_appearances"`
	Sort           string  `json:"sort" form:"sort"` // appearance_count.desc, total_duration.desc, first_appearance.asc, etc.
}

// New scene-based filtering for character scenes
type VideoSceneFilterAndPagination struct {
	models.BaseRequestParamsUri
	IncludeCharacters []uuid.UUID `json:"include_characters" form:"include_characters"`
	ExcludeCharacters []uuid.UUID `json:"exclude_characters" form:"exclude_characters"`
}

type VideoCharacterSummary struct {
	VideoID         uuid.UUID `json:"video_id"`
	CharacterID     uuid.UUID `json:"character_id"`
	CharacterName   string    `json:"character_name"`
	CharacterAvatar string    `json:"character_avatar"`
	DisplayName     string    `json:"display_name"`
	StartTime       float64   `json:"start_time"`
	EndTime         float64   `json:"end_time"`
}

type VideoCharacterListResponse struct {
	models.BaseListResponse
	Items []VideoCharacterSummary `json:"items"`
}

// New scene-based models
type VideoSceneCharacter struct {
	CharacterID     uuid.UUID `json:"character_id"`
	CharacterName   string    `json:"character_name"`
	CharacterAvatar string    `json:"character_avatar"`
	Confidence      float64   `json:"confidence"`
	StartTime       float64   `json:"start_time"`
	EndTime         float64   `json:"end_time"`
	StartFrame      int       `json:"start_frame"`
	EndFrame        int       `json:"end_frame"`
}

type VideoScene struct {
	VideoID            uuid.UUID             `json:"video_id"`
	SceneID            string                `json:"scene_id"`             // Generated based on time range
	StartTime          float64               `json:"start_time"`           // Scene start time
	EndTime            float64               `json:"end_time"`             // Scene end time
	Duration           float64               `json:"duration"`             // Scene duration
	StartFrame         int                   `json:"start_frame"`          // Scene start frame
	EndFrame           int                   `json:"end_frame"`            // Scene end frame
	CharacterCount     int                   `json:"character_count"`      // Number of characters in scene
	Characters         []VideoSceneCharacter `json:"characters"`           // Characters in this scene
	StartTimeFormatted string                `json:"start_time_formatted"` // HH:MM:SS
	EndTimeFormatted   string                `json:"end_time_formatted"`   // HH:MM:SS
}

type VideoSceneListResponse struct {
	models.BaseListResponse
	Items []VideoScene `json:"items"`
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

type VideoCharacterSummaryFilter struct {
	CharacterName  string  `json:"character_name"`
	MinConfidence  float64 `json:"min_confidence"`
	MinAppearances int     `json:"min_appearances"`
	Sort           string  `json:"sort"`
	Limit          int     `json:"limit"`
	Offset         int     `json:"offset"`
}
