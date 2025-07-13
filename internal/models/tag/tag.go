package tag

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

type TagPosition struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string `gorm:"type:text;not null" json:"title"`
	Position    string `gorm:"type:text;not null;unique" json:"position"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	models.Base
}

func (TagPosition) TableName() string {
	return common.POSTGRES_TABLE_NAME_TAG_POSITIONS
}

type TagCategory struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"type:text;not null;unique" json:"name"`
	Code             string    `gorm:"type:text;not null;unique" json:"code"`
	Description      string    `gorm:"type:text" json:"description"`
	Color            string    `gorm:"type:text;default:'#007bff'" json:"color"`
	Icon             string    `gorm:"type:text" json:"icon"`
	Priority         int       `gorm:"default:0" json:"priority"`
	IsShown          bool      `gorm:"default:true" json:"is_shown"`
	IsSystemCategory bool      `gorm:"default:false" json:"is_system_category"`
	FilterType       string    `gorm:"type:text;default:'single'" json:"filter_type"`
	CreatedBy        uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy        uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	models.Base
}

func (TagCategory) TableName() string {
	return common.POSTGRES_TABLE_NAME_TAG_CATEGORIES
}

type Tag struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  int       `gorm:"not null;index" json:"category_id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Code        string    `gorm:"type:text;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"type:text" json:"color"`
	Icon        string    `gorm:"type:text" json:"icon"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsSystemTag bool      `gorm:"default:false" json:"is_system_tag"`
	CreatedBy   uuid.UUID `gorm:"type:uuid" json:"created_by"`
	UpdatedBy   uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	models.Base

	Category TagCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Tag) TableName() string {
	return common.POSTGRES_TABLE_NAME_TAGS
}

type TagPositionCategory struct {
	ID            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	TagPositionID int    `gorm:"not null;index" json:"tag_position_id"`
	TagCategoryID int    `gorm:"not null;index" json:"tag_category_id"`
	SortOrder     int    `gorm:"default:0" json:"sort_order"`
	IsVisible     bool   `gorm:"default:true" json:"is_visible"`
	DisplayStyle  string `gorm:"type:text;default:'checkbox'" json:"display_style"`
	models.Base

	TagPosition TagPosition `gorm:"foreignKey:TagPositionID" json:"tag_position,omitempty"`
	TagCategory TagCategory `gorm:"foreignKey:TagCategoryID" json:"tag_category,omitempty"`
}

func (TagPositionCategory) TableName() string {
	return common.POSTGRES_TABLE_NAME_TAG_POSITION_CATEGORIES
}

type TagFilterRequest struct {
	models.BaseRequestParamsUri
	PositionCode string `form:"position" json:"position"`
	CategoryCode string `form:"category" json:"category"`
	IsActive     *bool  `form:"is_active" json:"is_active"`
	IsSystemTag  *bool  `form:"is_system_tag" json:"is_system_tag"`
}

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
