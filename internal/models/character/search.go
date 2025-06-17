package character

import (
	"github.com/google/uuid"
)

// CharacterSearchRequest represents complex search criteria
type CharacterSearchRequest struct {
	// Video filtering
	VideoIDs   []uuid.UUID `json:"video_ids" form:"video_ids"`
	VideoTitle string      `json:"video_title" form:"video_title"`

	// Character filtering
	IncludeCharacters []uuid.UUID `json:"include_characters" form:"include_characters"` // Must have ALL these characters
	ExcludeCharacters []uuid.UUID `json:"exclude_characters" form:"exclude_characters"` // Must NOT have ANY of these
	CharacterNames    []string    `json:"character_names" form:"character_names"`

	// Time filtering
	TimeRanges  []TimeRange `json:"time_ranges" form:"time_ranges"`
	MinDuration float64     `json:"min_duration" form:"min_duration"`
	MaxDuration float64     `json:"max_duration" form:"max_duration"`

	// Advanced filtering
	MinConfidence     float64 `json:"min_confidence" form:"min_confidence"`
	GroupByScenes     bool    `json:"group_by_scenes" form:"group_by_scenes"`
	SceneGapThreshold int     `json:"scene_gap_threshold" form:"scene_gap_threshold"` // seconds

	// Response options
	IncludeAppearances bool   `json:"include_appearances" form:"include_appearances"`
	IncludeScenes      bool   `json:"include_scenes" form:"include_scenes"`
	SortBy             string `json:"sort_by" form:"sort_by"`       // time, duration, character_count
	SortOrder          string `json:"sort_order" form:"sort_order"` // asc, desc

	// Pagination
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

// TimeRange represents a time range filter
type TimeRange struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// CharacterSearchResponse represents the search results
type CharacterSearchResponse struct {
	Videos        []VideoSearchResult `json:"videos"`
	Scenes        []SceneSegment      `json:"scenes,omitempty"`
	TotalCount    int                 `json:"total_count"`
	CurrentPage   int                 `json:"current_page"`
	PageSize      int                 `json:"page_size"`
	TotalPages    int                 `json:"total_pages"`
	SearchSummary SearchSummary       `json:"search_summary"`
}

// VideoSearchResult represents a video in search results
type VideoSearchResult struct {
	VideoID       uuid.UUID               `json:"video_id"`
	VideoTitle    string                  `json:"video_title"`
	VideoDuration int                     `json:"video_duration"`
	ThumbnailURL  string                  `json:"thumbnail_url"`
	Characters    []VideoCharacterSummary `json:"characters"`
	Appearances   []CharacterAppearance   `json:"appearances,omitempty"`
	SceneSegments []SceneSegment          `json:"scene_segments,omitempty"`
	MatchScore    float64                 `json:"match_score"` // Relevance score
}

// SearchSummary provides overall search statistics
type SearchSummary struct {
	TotalVideos      int             `json:"total_videos"`
	TotalCharacters  int             `json:"total_characters"`
	TotalAppearances int             `json:"total_appearances"`
	TotalDuration    float64         `json:"total_duration"`
	CharacterStats   []CharacterStat `json:"character_stats"`
	ExecutionTime    int             `json:"execution_time_ms"`
}

// CharacterStat provides statistics for individual characters
type CharacterStat struct {
	CharacterID     uuid.UUID `json:"character_id"`
	CharacterName   string    `json:"character_name"`
	VideoCount      int       `json:"video_count"`
	AppearanceCount int       `json:"appearance_count"`
	TotalDuration   float64   `json:"total_duration"`
}

// QuickSearchRequest for simple character search
type QuickSearchRequest struct {
	Query    string      `json:"query" form:"query"`
	VideoIDs []uuid.UUID `json:"video_ids" form:"video_ids"`
	Limit    int         `json:"limit" form:"limit"`
}

// QuickSearchResponse for simple search results
type QuickSearchResponse struct {
	Results []QuickSearchResult `json:"results"`
	Total   int                 `json:"total"`
}

// QuickSearchResult for individual quick search result
type QuickSearchResult struct {
	VideoID       uuid.UUID `json:"video_id"`
	VideoTitle    string    `json:"video_title"`
	CharacterName string    `json:"character_name"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	MatchScore    float64   `json:"match_score"`
}
