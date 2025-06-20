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
	GetVideoScenesWithCharacters(videoID string, queryParams characterModel.VideoSceneFilterAndPagination) (*characterModel.VideoSceneListResponse, error)
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

func (s *characterService) GetVideoScenesWithCharacters(videoID string, queryParams characterModel.VideoSceneFilterAndPagination) (*characterModel.VideoSceneListResponse, error) {
	uuidID, err := uuid.Parse(videoID)
	if err != nil {
		return nil, common.ErrInvalidUUID
	}

	queryParams.VerifyPaging()

	fmt.Printf("[DEBUG] GetVideoScenesWithCharacters - VideoID: %s\n", videoID)
	fmt.Printf("[DEBUG] Include Characters: %v\n", queryParams.IncludeCharacters)
	fmt.Printf("[DEBUG] Exclude Characters: %v\n", queryParams.ExcludeCharacters)
	fmt.Printf("[DEBUG] Page: %d, PageSize: %d\n", queryParams.Page, queryParams.PageSize)

	timeSegments, err := s.appearanceRepo.FindTimeSegmentsWithCharacters(s.sc.Ctx(), uuidID, queryParams.IncludeCharacters, queryParams.ExcludeCharacters)
	if err != nil {
		fmt.Printf("[DEBUG] Error in repository time segment finding: %v\n", err)
		return nil, err
	}

	fmt.Printf("[DEBUG] Repository returned time segments: %d\n", len(timeSegments))

	scenes := s.mapTimeSegmentsToVideoScenes(uuidID, timeSegments)

	fmt.Printf("[DEBUG] Mapped scenes: %d\n", len(scenes))

	total := len(scenes)
	offset := (queryParams.Page - 1) * queryParams.PageSize
	limit := queryParams.PageSize

	fmt.Printf("[DEBUG] Pagination - Total: %d, Offset: %d, Limit: %d\n", total, offset, limit)

	if offset >= total {
		scenes = []characterModel.VideoScene{}
	} else if offset+limit > total {
		scenes = scenes[offset:]
	} else {
		scenes = scenes[offset : offset+limit]
	}

	response := &characterModel.VideoSceneListResponse{
		BaseListResponse: models.BaseListResponse{
			Total:    total,
			Page:     queryParams.Page,
			PageSize: queryParams.PageSize,
		},
		Items: scenes,
	}

	fmt.Printf("[DEBUG] Final response - Items: %d\n", len(response.Items))
	return response, nil
}

func (s *characterService) mapTimeSegmentsToVideoScenes(videoID uuid.UUID, timeSegments []characterModel.TimeSegmentResult) []characterModel.VideoScene {
	var scenes []characterModel.VideoScene
	sceneIndex := 1

	for _, segment := range timeSegments {
		var sceneCharacters []characterModel.VideoSceneCharacter
		for _, char := range segment.Characters {
			sceneCharacters = append(sceneCharacters, characterModel.VideoSceneCharacter{
				CharacterID:     char.CharacterID,
				CharacterName:   char.CharacterName,
				CharacterAvatar: char.CharacterAvatar,
				Confidence:      char.Confidence,
				StartTime:       segment.StartTime,
				EndTime:         segment.EndTime,
			})
		}

		scene := characterModel.VideoScene{
			VideoID:            videoID,
			SceneID:            fmt.Sprintf("segment_%d_%.1f_%.1f", sceneIndex, segment.StartTime, segment.EndTime),
			StartTime:          segment.StartTime,
			EndTime:            segment.EndTime,
			Duration:           segment.Duration,
			CharacterCount:     len(sceneCharacters),
			Characters:         sceneCharacters,
			StartTimeFormatted: formatSecondsToTime(segment.StartTime),
			EndTimeFormatted:   formatSecondsToTime(segment.EndTime),
		}

		scenes = append(scenes, scene)
		fmt.Printf("[DEBUG] Service mapped scene %d: %.1f-%.1f (%.1fs) with %d characters\n",
			sceneIndex, segment.StartTime, segment.EndTime, segment.Duration, len(sceneCharacters))

		sceneIndex++
	}

	return scenes
}

func formatSecondsToTime(seconds float64) string {
	totalSeconds := int(seconds)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
