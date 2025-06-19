package video

import (
	"smart-scene-app-api/common"
	models "smart-scene-app-api/internal/models"

	"github.com/google/uuid"
)

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
	Items []VideoListingResponse `json:"items"`
}
