package video

import (
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

// VideoListingRequest - Request cho API listing videos
type VideoListingRequest struct {
	models.BaseRequestParamsUri

	// Search parameters
	Search         string   `form:"search" json:"search"`                   // Search theo title
	CharacterNames []string `form:"character_names" json:"character_names"` // Search theo character names
	CharacterIDs   []string `form:"character_ids" json:"character_ids"`     // Search theo character IDs

	// Filter parameters
	Status   string   `form:"status" json:"status"`       // pending, processing, completed, failed
	TagIDs   []int    `form:"tag_ids" json:"tag_ids"`     // Filter theo tag IDs
	TagCodes []string `form:"tag_codes" json:"tag_codes"` // Filter theo tag codes (male, female, etc.)

	// Date range filters
	DateFrom string `form:"date_from" json:"date_from"` // YYYY-MM-DD
	DateTo   string `form:"date_to" json:"date_to"`     // YYYY-MM-DD

	// Video properties
	MinDuration   int   `form:"min_duration" json:"min_duration"`     // Minimum duration in seconds
	MaxDuration   int   `form:"max_duration" json:"max_duration"`     // Maximum duration in seconds
	HasCharacters *bool `form:"has_characters" json:"has_characters"` // Filter videos with/without characters

	// User filter
	CreatedBy uuid.UUID `form:"created_by" json:"created_by"`
}

// VideoListingResponse - Response cho từng video trong list
type VideoListingResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	ThumbnailURL   string    `json:"thumbnail_url"`
	Duration       int       `json:"duration"`
	CharacterCount int       `json:"character_count"`
	Status         string    `json:"status"`
	CreatedAt      string    `json:"created_at"` // ISO format
	UpdatedAt      string    `json:"updated_at"` // ISO format

	// Tags hiển thị trên thumbnail (Tom, Jerry, Mickey, +10)
	Tags             []VideoTagInfo `json:"tags"`
	VisibleTagsCount int            `json:"visible_tags_count"` // Số lượng tags hiển thị (3-4 tags)
	TotalTagsCount   int            `json:"total_tags_count"`   // Tổng số tags (+10)

	// Video properties
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format   string `json:"format"`
	FilePath string `json:"file_path"`
}

// VideoTagInfo - Thông tin tag hiển thị trên video thumbnail
type VideoTagInfo struct {
	TagID        int    `json:"tag_id"`
	TagName      string `json:"tag_name"`
	TagCode      string `json:"tag_code"`
	TagColor     string `json:"tag_color"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Priority     int    `json:"priority"` // Priority để sort tags
}

// VideoListResponse - Response chính cho API listing
type VideoListResponse struct {
	models.BaseListResponse
	Items   []VideoListingResponse `json:"items"`
	Filters VideoListingFilters    `json:"filters"` // Available filters for frontend
}

// VideoListingFilters - Available filters cho frontend
type VideoListingFilters struct {
	Statuses       []FilterOption        `json:"statuses"`
	DurationRanges []DurationRangeOption `json:"duration_ranges"`
	Tags           []TagFilterGroup      `json:"tags"` // Grouped tags by category
}

// FilterOption - Generic filter option
type FilterOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"` // Số lượng videos có option này
}

// DurationRangeOption - Duration range filter
type DurationRangeOption struct {
	MinDuration int    `json:"min_duration"`
	MaxDuration int    `json:"max_duration"`
	Label       string `json:"label"` // "0-30s", "30s-1m", "1m-5m", etc.
	Count       int    `json:"count"`
}

// TagFilterGroup - Tags grouped by category for sidebar filter
type TagFilterGroup struct {
	CategoryID   int               `json:"category_id"`
	CategoryName string            `json:"category_name"`
	CategoryCode string            `json:"category_code"`
	FilterType   string            `json:"filter_type"`   // "single", "multiple"
	DisplayStyle string            `json:"display_style"` // "checkbox", "radio", "dropdown"
	Tags         []TagFilterOption `json:"tags"`
}

// TagFilterOption - Individual tag option in filter
type TagFilterOption struct {
	TagID      int    `json:"tag_id"`
	TagName    string `json:"tag_name"`
	TagCode    string `json:"tag_code"`
	TagColor   string `json:"tag_color"`
	Count      int    `json:"count"`       // Số videos có tag này
	IsSelected bool   `json:"is_selected"` // Based on request params
}

// VideoSearchSuggestion - Suggestions cho search autocomplete
type VideoSearchSuggestion struct {
	Type         string `json:"type"` // "video", "character"
	ID           string `json:"id"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"` // Character name for videos, video count for characters
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

// VideoSearchSuggestionResponse - Response cho search suggestions
type VideoSearchSuggestionResponse struct {
	Videos     []VideoSearchSuggestion `json:"videos"`
	Characters []VideoSearchSuggestion `json:"characters"`
	Total      int                     `json:"total"`
}
