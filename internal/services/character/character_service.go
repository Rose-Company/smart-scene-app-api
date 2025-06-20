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

	var filters []repositories.Clause

	filters = append(filters, func(tx *gorm.DB) {
		tx.Where("video_id = ?", uuidID)
	})

	if queryParams.CharacterName != "" {
		var characterIDs []uuid.UUID
		characters, err := s.appearanceRepo.List(s.sc.Ctx(), models.QueryParams{}, func(tx *gorm.DB) {
		}, func(tx *gorm.DB) {
			tx.Where("name ILIKE ?", "%"+queryParams.CharacterName+"%").Select("id")
		})
		if err != nil {
			return nil, err
		}

		for _, char := range characters {
			characterIDs = append(characterIDs, char.CharacterID)
		}

		if len(characterIDs) > 0 {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where("character_id IN ?", characterIDs)
			})
		} else {
			return &characterModel.VideoCharacterListResponse{
				BaseListResponse: models.BaseListResponse{
					Total:    0,
					Page:     queryParams.Page,
					PageSize: queryParams.PageSize,
					Items:    []characterModel.VideoCharacterSummary{},
				},
			}, nil
		}
	}

	combinedFilter := func(tx *gorm.DB) {
		for _, filter := range filters {
			filter(tx)
		}
	}

	appearances, err := s.appearanceRepo.List(s.sc.Ctx(), models.QueryParams{}, combinedFilter)
	if err != nil {
		return nil, err
	}

	if len(appearances) == 0 {
		return &characterModel.VideoCharacterListResponse{
			BaseListResponse: models.BaseListResponse{
				Total:    0,
				Page:     queryParams.Page,
				PageSize: queryParams.PageSize,
				Items:    []characterModel.VideoCharacterSummary{},
			},
		}, nil
	}

	characterIDSet := make(map[uuid.UUID]bool)
	for _, appearance := range appearances {
		if appearance != nil {
			characterIDSet[appearance.CharacterID] = true
		}
	}

	characterIDs := make([]uuid.UUID, 0, len(characterIDSet))
	for id := range characterIDSet {
		characterIDs = append(characterIDs, id)
	}

	characterMap := make(map[uuid.UUID]*characterModel.Character)
	if len(characterIDs) > 0 {
		characters, err := s.characterRepo.List(s.sc.Ctx(), models.QueryParams{}, func(tx *gorm.DB) {
		}, func(tx *gorm.DB) {
			tx.Where("id IN ? AND is_active = ?", characterIDs, true)
		})
		if err != nil {
			return nil, err
		}

		for _, char := range characters {
			if char != nil {
				characterMap[char.ID] = char
			}
		}
	}

	characters := make([]characterModel.Character, 0, len(characterMap))
	for _, character := range characterMap {
		if character != nil {
			characters = append(characters, *character)
		}
	}

	total := len(characters)
	offset := (queryParams.Page - 1) * queryParams.PageSize
	limit := queryParams.PageSize

	if offset >= total {
		characters = []characterModel.Character{}
	} else if offset+limit > total {
		characters = characters[offset:]
	} else {
		characters = characters[offset : offset+limit]
	}

	items := make([]characterModel.VideoCharacterSummary, 0, len(characters))

	for _, character := range characters {
		item := characterModel.VideoCharacterSummary{
			VideoID:         uuidID,
			CharacterID:     character.ID,
			CharacterName:   character.Name,
			CharacterAvatar: character.Avatar,
		}
		items = append(items, item)
	}

	response := &characterModel.VideoCharacterListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    len(items),
			Page:     queryParams.Page,
			PageSize: queryParams.PageSize,
		},
		Items: items,
	}

	return response, nil
}

func formatSecondsToTime(seconds float64) string {
	totalSeconds := int(seconds)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
