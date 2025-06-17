package tag

import (
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

// TagPosition - Vị trí hiển thị tags (sidebar, thumbnail, etc.)
type TagPosition struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"type:text;not null" json:"title"`           // "Video Theme", "People"
	Position    string `gorm:"type:text;not null;unique" json:"position"` // "video_theme", "people"
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	models.Base
}

// TagCategory - Loại tags (Gender, Age, Character Type, etc.)
type TagCategory struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"type:text;not null;unique" json:"name"` // "Gender", "Age Range"
	Code             string    `gorm:"type:text;not null;unique" json:"code"` // "gender", "age_range"
	Description      string    `gorm:"type:text" json:"description"`
	Color            string    `gorm:"type:text;default:'#007bff'" json:"color"`
	Icon             string    `gorm:"type:text" json:"icon"`
	Priority         int       `gorm:"default:0" json:"priority"`
	IsShown          bool      `gorm:"default:true" json:"is_shown"`
	IsSystemCategory bool      `gorm:"default:false" json:"is_system_category"`
	FilterType       string    `gorm:"type:text;default:'single'" json:"filter_type"` // "single", "multiple", "range"
	CreatedBy        uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy        uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	models.Base
}

// Tag - Tags cụ thể (Male, Female, Child, etc.)
type Tag struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  int       `gorm:"not null;index" json:"category_id"`
	Name        string    `gorm:"type:text;not null" json:"name"` // "Male", "Female"
	Code        string    `gorm:"type:text;not null" json:"code"` // "male", "female"
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"type:text" json:"color"` // Inherit từ category nếu NULL
	Icon        string    `gorm:"type:text" json:"icon"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsSystemTag bool      `gorm:"default:false" json:"is_system_tag"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy   uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	models.Base

	// Relations
	Category TagCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TagPositionCategory - Mapping categories với positions
type TagPositionCategory struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	TagPositionID int    `gorm:"not null;index" json:"tag_position_id"`
	TagCategoryID int    `gorm:"not null;index" json:"tag_category_id"`
	SortOrder     int    `gorm:"default:0" json:"sort_order"`
	IsVisible     bool   `gorm:"default:true" json:"is_visible"`
	DisplayStyle  string `gorm:"type:text;default:'checkbox'" json:"display_style"` // "dropdown", "checkbox", "radio", "chips"
	models.Base

	// Relations
	TagPosition TagPosition `gorm:"foreignKey:TagPositionID" json:"tag_position,omitempty"`
	TagCategory TagCategory `gorm:"foreignKey:TagCategoryID" json:"tag_category,omitempty"`
}

// Request/Response Models
type TagFilterRequest struct {
	models.BaseRequestParamsUri
	PositionCode string `form:"position" json:"position"` // "people", "video_theme"
	CategoryCode string `form:"category" json:"category"` // "gender", "age_range"
	IsActive     *bool  `form:"is_active" json:"is_active"`
	IsSystemTag  *bool  `form:"is_system_tag" json:"is_system_tag"`
}

// TagHierarchyResponse - Response cho hierarchical tags
type TagHierarchyResponse struct {
	PositionID    int                   `json:"position_id"`
	PositionTitle string                `json:"position_title"`
	PositionCode  string                `json:"position_code"`
	Categories    []TagCategoryResponse `json:"categories"`
}

type TagCategoryResponse struct {
	CategoryID   int           `json:"category_id"`
	CategoryName string        `json:"category_name"`
	CategoryCode string        `json:"category_code"`
	Color        string        `json:"color"`
	FilterType   string        `json:"filter_type"`
	DisplayStyle string        `json:"display_style"`
	Tags         []TagResponse `json:"tags"`
}

type TagResponse struct {
	TagID      int    `json:"tag_id"`
	TagName    string `json:"tag_name"`
	TagCode    string `json:"tag_code"`
	Color      string `json:"color"`
	UsageCount int    `json:"usage_count"`
	IsActive   bool   `json:"is_active"`
}

type TagListResponse struct {
	models.BaseListResponse
	Items []TagHierarchyResponse `json:"items"`
}
