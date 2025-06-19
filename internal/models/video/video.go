package video

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID                   uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
	Status               string      `json:"status" gorm:"type:text;not null;check:status IN ('pending', 'processing', 'completed', 'failed')"`
	CreatedBy            uuid.UUID   `json:"created_by" gorm:"type:uuid;not null;default:gen_random_uuid()"`
	UpdatedBy            uuid.UUID   `json:"updated_by" gorm:"type:uuid;not null;default:gen_random_uuid()"`
	Title                string      `json:"title" gorm:"type:text;not null"`
	FilePath             string      `json:"file_path" gorm:"type:text;not null"`
	Duration             int         `json:"duration" gorm:"type:int;not null"`
	Metadata             common.JSON `json:"metadata" gorm:"type:jsonb"`
	ThumbnailURL         string      `json:"thumbnail_url" gorm:"type:text"`
	HasCharacterAnalysis bool        `json:"has_character_analysis" gorm:"default:false"`
	CharacterCount       int         `json:"character_count" gorm:"type:int;default:0"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoFilterAndPagination struct {
	models.BaseRequestParamsUri
	Title     string    `json:"title" form:"title"`
	Status    string    `json:"status" form:"status"`
	CreatedBy uuid.UUID `json:"created_by" form:"created_by"`
	TagIDs    []int     `json:"tag_ids" form:"tag_ids"`
	TagCodes  []string  `json:"tag_codes" form:"tag_codes"`
}
