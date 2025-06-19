package character

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	characterModel "smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"
	characterRepo "smart-scene-app-api/internal/repositories/character"
	"smart-scene-app-api/server"

	"fmt"

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
	uuidID, err := uuid.Parse(videoID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	queryParams.VerifyPaging()

	limit := queryParams.PageSize
	offset := (queryParams.Page - 1) * queryParams.PageSize

	var filters []repositories.Clause

	filters = append(filters, func(tx *gorm.DB) {
		tx.Where("video_id = ?", uuidID)
	})

	if queryParams.CharacterName != "" {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Joins("JOIN characters c ON character_appearances.character_id = c.id").
				Where("c.name ILIKE ?", "%"+queryParams.CharacterName+"%")
		})
	}

	if queryParams.MinConfidence > 0 {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("confidence >= ?", queryParams.MinConfidence)
		})
	}

	total, err := s.appearanceRepo.Count(s.sc.Ctx(), models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	response := &characterModel.VideoCharacterListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    int(total),
			Page:     queryParams.Page,
			PageSize: queryParams.PageSize,
			Items:    []characterModel.VideoCharacterSummary{},
		},
	}

	if total == 0 {
		return response, nil
	}

	sort := queryParams.Sort
	if sort == "" {
		sort = "start_time.asc"
	}

	repoQueryParams := models.QueryParams{
		Limit:  limit,
		Offset: offset,
		QuerySort: models.QuerySort{
			Origin: sort,
		},
	}

	preloadFunc := func(tx *gorm.DB) {
		tx.Preload("Character")
	}

	combinedFilter := func(tx *gorm.DB) {
		for _, filter := range filters {
			filter(tx)
		}
	}

	appearances, err := s.appearanceRepo.List(s.sc.Ctx(), repoQueryParams, preloadFunc, combinedFilter)
	if err != nil {
		return nil, err
	}

	characterSummaryMap := make(map[uuid.UUID]*characterModel.VideoCharacterSummary)
	confidenceSum := make(map[uuid.UUID]float64)
	firstAppearanceTime := make(map[uuid.UUID]float64)
	lastAppearanceTime := make(map[uuid.UUID]float64)

	for _, appearance := range appearances {
		if appearance != nil && appearance.Character != nil {
			charID := appearance.CharacterID

			if summary, exists := characterSummaryMap[charID]; exists {
				summary.AppearanceCount++
				summary.TotalDuration += (appearance.EndTime - appearance.StartTime)
				confidenceSum[charID] += appearance.Confidence

				if appearance.StartTime < firstAppearanceTime[charID] {
					firstAppearanceTime[charID] = appearance.StartTime
					summary.FirstAppearance = formatSecondsToTime(appearance.StartTime)
					summary.FirstAppearanceFrame = appearance.StartFrame
				}

				if appearance.EndTime > lastAppearanceTime[charID] {
					lastAppearanceTime[charID] = appearance.EndTime
					summary.LastAppearance = formatSecondsToTime(appearance.EndTime)
					summary.LastAppearanceFrame = appearance.EndFrame
				}

				summary.AvgConfidence = confidenceSum[charID] / float64(summary.AppearanceCount)
			} else {
				firstAppearanceTime[charID] = appearance.StartTime
				lastAppearanceTime[charID] = appearance.EndTime
				confidenceSum[charID] = appearance.Confidence

				summary := &characterModel.VideoCharacterSummary{
					VideoID:              appearance.VideoID,
					CharacterID:          charID,
					CharacterName:        appearance.Character.Name,
					CharacterAvatar:      appearance.Character.Avatar,
					DisplayName:          appearance.Character.Name,
					AppearanceCount:      1,
					TotalDuration:        appearance.EndTime - appearance.StartTime,
					FirstAppearance:      formatSecondsToTime(appearance.StartTime),
					LastAppearance:       formatSecondsToTime(appearance.EndTime),
					AvgConfidence:        appearance.Confidence,
					FirstAppearanceFrame: appearance.StartFrame,
					LastAppearanceFrame:  appearance.EndFrame,
				}
				characterSummaryMap[charID] = summary
			}
		}
	}

	items := make([]characterModel.VideoCharacterSummary, 0, len(characterSummaryMap))
	for _, summary := range characterSummaryMap {
		items = append(items, *summary)
	}

	response.Items = items
	return response, nil
}

func formatSecondsToTime(seconds float64) string {
	totalSeconds := int(seconds)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
