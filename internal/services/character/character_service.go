package character

import (
	"fmt"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	characterModel "smart-scene-app-api/internal/models/character"
	characterRepo "smart-scene-app-api/internal/repositories/character"
	"smart-scene-app-api/server"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetCharactersByVideoID(videoID string, queryParams characterModel.VideoCharacterFilterAndPagination) (*characterModel.VideoCharacterListResponse, error)
}

type characterService struct {
	sc             server.ServerContext
	characterRepo  characterRepo.Repository
	appearanceRepo characterRepo.AppearanceRepository
}

func NewCharacterService(sc server.ServerContext) Service {
	return &characterService{
		sc:             sc,
		characterRepo:  characterRepo.NewRepository(sc.DB()),
		appearanceRepo: characterRepo.NewAppearanceRepository(sc.DB()),
	}
}

func (s *characterService) GetCharactersByVideoID(videoID string, queryParams characterModel.VideoCharacterFilterAndPagination) (*characterModel.VideoCharacterListResponse, error) {
	// Parse video ID
	uuidVideoID, err := uuid.Parse(videoID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	// Verify paging parameters
	queryParams.VerifyPaging()

	limit := queryParams.PageSize
	offset := (queryParams.Page - 1) * queryParams.PageSize

	// Get character summary for the video with filters
	var filters []func(*gorm.DB)

	// Add video filter
	filters = append(filters, func(tx *gorm.DB) {
		tx.Where("ca.video_id = ?", uuidVideoID)
	})

	// Add character name filter if provided
	if queryParams.CharacterName != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("c.name ILIKE ?", "%"+queryParams.CharacterName+"%")
		})
	}

	// Add confidence filter if provided
	if queryParams.MinConfidence > 0 {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Having("AVG(ca.confidence) >= ?", queryParams.MinConfidence)
		})
	}

	// Add appearance count filter if provided
	if queryParams.MinAppearances > 0 {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Having("COUNT(ca.id) >= ?", queryParams.MinAppearances)
		})
	}

	// Count total characters for pagination
	var total int64
	countQuery := s.sc.DB().
		Table("character_appearances ca").
		Select("COUNT(DISTINCT ca.character_id)").
		Joins("JOIN characters c ON ca.character_id = c.id")

	for _, filter := range filters {
		filter(countQuery)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Prepare response
	response := &characterModel.VideoCharacterListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    int(total),
			Page:     queryParams.Page,
			PageSize: queryParams.PageSize,
			Items:    []characterModel.VideoCharacterSummary{},
		},
	}

	// If no characters found, return empty response
	if total == 0 {
		return response, nil
	}

	// Build main query with aggregations
	query := s.sc.DB().
		Table("character_appearances ca").
		Select(`
			ca.video_id,
			ca.character_id,
			c.name as character_name,
			c.avatar as character_avatar,
			c.display_name,
			COUNT(ca.id) as appearance_count,
			SUM(CASE WHEN ca.end_time > ca.start_time THEN ca.end_time - ca.start_time ELSE 0 END) as total_duration,
			MIN(ca.start_time) as first_appearance_time,
			MAX(ca.end_time) as last_appearance_time,
			AVG(ca.confidence) as avg_confidence,
			MIN(ca.start_frame) as first_appearance_frame,
			MAX(ca.end_frame) as last_appearance_frame
		`).
		Joins("JOIN characters c ON ca.character_id = c.id").
		Group("ca.video_id, ca.character_id, c.name, c.avatar, c.display_name")

	// Apply filters
	for _, filter := range filters {
		filter(query)
	}

	// Apply sorting
	sort := queryParams.Sort
	if sort == "" {
		sort = "appearance_count.desc" // Default sort by appearance count descending
	}

	switch sort {
	case "appearance_count.asc":
		query = query.Order("appearance_count ASC")
	case "appearance_count.desc":
		query = query.Order("appearance_count DESC")
	case "total_duration.asc":
		query = query.Order("total_duration ASC")
	case "total_duration.desc":
		query = query.Order("total_duration DESC")
	case "first_appearance.asc":
		query = query.Order("first_appearance_time ASC")
	case "first_appearance.desc":
		query = query.Order("first_appearance_time DESC")
	case "character_name.asc":
		query = query.Order("character_name ASC")
	case "character_name.desc":
		query = query.Order("character_name DESC")
	case "confidence.asc":
		query = query.Order("avg_confidence ASC")
	case "confidence.desc":
		query = query.Order("avg_confidence DESC")
	default:
		query = query.Order("appearance_count DESC")
	}

	// Apply pagination
	query = query.Limit(limit).Offset(offset)

	// Execute query
	var summaries []struct {
		VideoID              uuid.UUID `json:"video_id"`
		CharacterID          uuid.UUID `json:"character_id"`
		CharacterName        string    `json:"character_name"`
		CharacterAvatar      string    `json:"character_avatar"`
		DisplayName          string    `json:"display_name"`
		AppearanceCount      int       `json:"appearance_count"`
		TotalDuration        float64   `json:"total_duration"`
		FirstAppearanceTime  float64   `json:"first_appearance_time"`
		LastAppearanceTime   float64   `json:"last_appearance_time"`
		AvgConfidence        float64   `json:"avg_confidence"`
		FirstAppearanceFrame int       `json:"first_appearance_frame"`
		LastAppearanceFrame  int       `json:"last_appearance_frame"`
	}

	if err := query.Scan(&summaries).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	items := make([]characterModel.VideoCharacterSummary, 0, len(summaries))
	for _, summary := range summaries {
		// Convert time to HH:MM:SS format
		firstAppearance := formatSecondsToTime(summary.FirstAppearanceTime)
		lastAppearance := formatSecondsToTime(summary.LastAppearanceTime)

		item := characterModel.VideoCharacterSummary{
			VideoID:              summary.VideoID,
			CharacterID:          summary.CharacterID,
			CharacterName:        summary.CharacterName,
			CharacterAvatar:      summary.CharacterAvatar,
			DisplayName:          summary.DisplayName,
			AppearanceCount:      summary.AppearanceCount,
			TotalDuration:        summary.TotalDuration,
			FirstAppearance:      firstAppearance,
			LastAppearance:       lastAppearance,
			AvgConfidence:        summary.AvgConfidence,
			FirstAppearanceFrame: summary.FirstAppearanceFrame,
			LastAppearanceFrame:  summary.LastAppearanceFrame,
		}
		items = append(items, item)
	}

	response.Items = items
	return response, nil
}

// Helper function to format seconds to HH:MM:SS
func formatSecondsToTime(seconds float64) string {
	totalSeconds := int(seconds)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
