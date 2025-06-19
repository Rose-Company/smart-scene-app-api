package video

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

type VideoListingRequest struct {
	models.BaseRequestParamsUri
	Search         string    `form:"search" json:"search"`
	CharacterNames []string  `form:"character_names" json:"character_names"`
	CharacterIDs   []string  `form:"character_ids" json:"character_ids"`
	Status         string    `form:"status" json:"status"`
	TagIDs         []int     `form:"tag_ids" json:"tag_ids"`
	TagCodes       []string  `form:"tag_codes" json:"tag_codes"`
	DateFrom       string    `form:"date_from" json:"date_from"`
	DateTo         string    `form:"date_to" json:"date_to"`
	MinDuration    int       `form:"min_duration" json:"min_duration"`
	MaxDuration    int       `form:"max_duration" json:"max_duration"`
	HasCharacters  *bool     `form:"has_characters" json:"has_characters"`
	CreatedBy      uuid.UUID `form:"created_by" json:"created_by"`
}

type VideoListingResponse struct {
	ID               uuid.UUID      `json:"id"`
	Title            string         `json:"title"`
	ThumbnailURL     string         `json:"thumbnail_url"`
	Duration         int            `json:"duration"`
	CharacterCount   int            `json:"character_count"`
	Status           string         `json:"status"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	FilePath         string         `json:"file_path"`
	Tags             []VideoTagInfo `json:"tags"`
	VisibleTagsCount int            `json:"visible_tags_count"`
	TotalTagsCount   int            `json:"total_tags_count"`
	Metadata         common.JSON    `json:"metadata"`
}

type VideoTagInfo struct {
	TagID        int    `json:"tag_id"`
	TagName      string `json:"tag_name"`
	TagCode      string `json:"tag_code"`
	TagColor     string `json:"tag_color"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Priority     int    `json:"priority"`
}

type VideoListResponse struct {
	models.BaseListResponse
	Items   []VideoListingResponse `json:"items"`
	Filters VideoListingFilters    `json:"filters"`
}

type VideoListingFilters struct {
	Statuses       []FilterOption        `json:"statuses"`
	DurationRanges []DurationRangeOption `json:"duration_ranges"`
	Tags           []TagFilterGroup      `json:"tags"`
}

type FilterOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type DurationRangeOption struct {
	MinDuration int    `json:"min_duration"`
	MaxDuration int    `json:"max_duration"`
	Label       string `json:"label"`
	Count       int    `json:"count"`
}

type TagFilterGroup struct {
	CategoryID   int               `json:"category_id"`
	CategoryName string            `json:"category_name"`
	CategoryCode string            `json:"category_code"`
	FilterType   string            `json:"filter_type"`
	DisplayStyle string            `json:"display_style"`
	Tags         []TagFilterOption `json:"tags"`
}

type TagFilterOption struct {
	TagID      int    `json:"tag_id"`
	TagName    string `json:"tag_name"`
	TagCode    string `json:"tag_code"`
	TagColor   string `json:"tag_color"`
	Count      int    `json:"count"`
	IsSelected bool   `json:"is_selected"`
}

type VideoSearchSuggestion struct {
	Type         string `json:"type"`
	ID           string `json:"id"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

type VideoSearchSuggestionResponse struct {
	Videos     []VideoSearchSuggestion `json:"videos"`
	Characters []VideoSearchSuggestion `json:"characters"`
	Total      int                     `json:"total"`
}
