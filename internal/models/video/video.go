package video

import (
	models "smart-scene-app-api/internal/models"
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID                   uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	Status               string                 `json:"status" gorm:"type:text;not null;check:status IN ('pending', 'processing', 'completed', 'failed')"`
	CreatedBy            uuid.UUID              `json:"created_by" gorm:"type:uuid;not null;default:gen_random_uuid()"`
	UpdatedBy            uuid.UUID              `json:"updated_by" gorm:"type:uuid;not null;default:gen_random_uuid()"`
	Title                string                 `json:"title" gorm:"type:text;not null"`
	FilePath             string                 `json:"file_path" gorm:"type:text;not null"`
	Duration             int                    `json:"duration" gorm:"type:int;not null"`
	Width                int                    `json:"width" gorm:"type:int"`
	Height               int                    `json:"height" gorm:"type:int"`
	Folder               string                 `json:"folder" gorm:"type:text"`
	Format               string                 `json:"format" gorm:"type:text"`
	Metadata             map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	ThumbnailURL         string                 `json:"thumbnail_url" gorm:"type:text"`
	FrameRate            float64                `json:"frame_rate" gorm:"type:decimal(8,4);default:24.0"`
	HasCharacterAnalysis bool                   `json:"has_character_analysis" gorm:"default:false"`
	CharacterCount       int                    `json:"character_count" gorm:"type:int;default:0"`
	AnalysisVersion      string                 `json:"analysis_version" gorm:"type:text"`
	Analysis             interface{}            `json:"analysis,omitempty" gorm:"foreignKey:VideoID"`
	CharacterAppearances interface{}            `json:"character_appearances,omitempty" gorm:"foreignKey:VideoID"`
}

func (Video) TableName() string {
	return "videos"
}

type VideoFilterAndPagination struct {
	Title       string             `json:"title" form:"title"`
	Status      string             `json:"status" form:"status"`
	CreatedBy   uuid.UUID          `json:"created_by" form:"created_by"`
	QueryParams models.QueryParams `json:"query_params" form:"query_params"`
}
