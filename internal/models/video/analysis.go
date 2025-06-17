package video

import (
	models "smart-scene-app-api/internal/models"
	"time"

	"github.com/google/uuid"
)

// VideoAnalysis represents the analysis status and results of a video
type VideoAnalysis struct {
	models.Base
	VideoID             uuid.UUID              `json:"video_id" gorm:"type:uuid;not null;index;unique"`
	Status              string                 `json:"status" gorm:"type:text;not null;check:status IN ('pending', 'processing', 'completed', 'failed');default:'pending'"`
	ProcessingType      string                 `json:"processing_type" gorm:"type:text;not null;default:'character_detection'"` // character_detection, scene_analysis, etc.
	Progress            int                    `json:"progress" gorm:"type:int;default:0;check:progress >= 0 AND progress <= 100"`
	ProcessingStarted   *time.Time             `json:"processing_started,omitempty" gorm:"type:timestamp"`
	ProcessingCompleted *time.Time             `json:"processing_completed,omitempty" gorm:"type:timestamp"`
	ProcessingDuration  int                    `json:"processing_duration" gorm:"type:int;default:0"` // in seconds
	CharacterCount      int                    `json:"character_count" gorm:"type:int;default:0"`
	AppearanceCount     int                    `json:"appearance_count" gorm:"type:int;default:0"`
	AnalysisResults     map[string]interface{} `json:"analysis_results" gorm:"type:jsonb"`
	ErrorMessage        string                 `json:"error_message" gorm:"type:text"`
	ProcessedBy         uuid.UUID              `json:"processed_by" gorm:"type:uuid"`

	// Configurations
	Config map[string]interface{} `json:"config" gorm:"type:jsonb"` // AI model configs, thresholds, etc.

	// Relations
	Video interface{} `json:"video,omitempty" gorm:"foreignKey:VideoID"`
}

// VideoAnalysisFilterAndPagination for filtering analysis results
type VideoAnalysisFilterAndPagination struct {
	VideoID        uuid.UUID          `json:"video_id" form:"video_id"`
	Status         string             `json:"status" form:"status"`
	ProcessingType string             `json:"processing_type" form:"processing_type"`
	MinProgress    int                `json:"min_progress" form:"min_progress"`
	ProcessedBy    uuid.UUID          `json:"processed_by" form:"processed_by"`
	QueryParams    models.QueryParams `json:"query_params" form:"query_params"`
}
