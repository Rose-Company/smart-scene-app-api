package character

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	characterModel "smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"
	characterRepo "smart-scene-app-api/internal/repositories/character"
	"smart-scene-app-api/server"

	"fmt"
	"sort"

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

	// Merge overlapping segments and find all characters in merged ranges
	scenes, err := s.mapTimeSegmentsToVideoScenesWithMerging(uuidID, timeSegments, queryParams.IncludeCharacters)
	if err != nil {
		fmt.Printf("[DEBUG] Error in mapping segments to scenes: %v\n", err)
		return nil, err
	}

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

// mapTimeSegmentsToVideoScenesWithMerging merges overlapping segments and finds all characters in merged ranges
func (s *characterService) mapTimeSegmentsToVideoScenesWithMerging(videoID uuid.UUID, timeSegments []characterModel.TimeSegmentResult, requiredCharacters []uuid.UUID) ([]characterModel.VideoScene, error) {
	if len(timeSegments) == 0 {
		return []characterModel.VideoScene{}, nil
	}

	// Step 1: Merge overlapping time ranges
	mergedRanges := s.mergeTimeRanges(timeSegments)
	fmt.Printf("[DEBUG] Merged %d segments into %d ranges\n", len(timeSegments), len(mergedRanges))

	// Step 2: For each merged range, find all characters that appear in that time range
	var scenes []characterModel.VideoScene
	sceneCounter := 1
	for i, timeRange := range mergedRanges {
		sceneCharacters, err := s.findCharactersInTimeRange(videoID, timeRange.StartTime, timeRange.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to find characters in time range %.1f-%.1f: %w", timeRange.StartTime, timeRange.EndTime, err)
		}

		// Only create scene if there are characters present and all required characters are included
		if len(sceneCharacters) > 0 && s.sceneContainsAllRequiredCharacters(sceneCharacters, requiredCharacters) {
			// Use the intersection time from characters if available
			sceneStartTime := timeRange.StartTime
			sceneEndTime := timeRange.EndTime

			if len(sceneCharacters) > 0 {
				// All characters should have the same intersection time after findCharactersInTimeRange
				sceneStartTime = sceneCharacters[0].StartTime
				sceneEndTime = sceneCharacters[0].EndTime
			}

			scene := characterModel.VideoScene{
				VideoID:            videoID,
				SceneID:            fmt.Sprintf("segment_%d_%.1f_%.1f", sceneCounter, timeRange.StartTime, timeRange.EndTime),
				StartTime:          sceneStartTime,
				EndTime:            sceneEndTime,
				Duration:           sceneEndTime - sceneStartTime,
				CharacterCount:     len(sceneCharacters),
				Characters:         sceneCharacters,
				StartTimeFormatted: formatSecondsToTime(sceneStartTime),
				EndTimeFormatted:   formatSecondsToTime(sceneEndTime),
			}

			scenes = append(scenes, scene)
			fmt.Printf("[DEBUG] Created merged scene %d: %.1f-%.1f (intersection: %.1f-%.1f) with %d characters\n", sceneCounter, timeRange.StartTime, timeRange.EndTime, sceneStartTime, sceneEndTime, len(sceneCharacters))
			sceneCounter++
		} else {
			fmt.Printf("[DEBUG] Skipped scene %d: %.1f-%.1f (no characters found)\n", i+1, timeRange.StartTime, timeRange.EndTime)
		}
	}

	return scenes, nil
}

// TimeRange represents a time range
type TimeRange struct {
	StartTime float64
	EndTime   float64
}

// mergeTimeRanges merges overlapping time ranges from segments
func (s *characterService) mergeTimeRanges(timeSegments []characterModel.TimeSegmentResult) []TimeRange {
	if len(timeSegments) == 0 {
		return []TimeRange{}
	}

	// Convert to time ranges and sort by start time
	var ranges []TimeRange
	for _, segment := range timeSegments {
		ranges = append(ranges, TimeRange{
			StartTime: segment.StartTime,
			EndTime:   segment.EndTime,
		})
	}

	// Sort by start time
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].StartTime < ranges[j].StartTime
	})

	// Merge overlapping ranges
	var merged []TimeRange
	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		next := ranges[i]

		// If current and next overlap, merge them
		if current.EndTime >= next.StartTime {
			// Extend current range to include next
			if next.EndTime > current.EndTime {
				current.EndTime = next.EndTime
			}
			fmt.Printf("[DEBUG] Merged ranges: %.1f-%.1f + %.1f-%.1f = %.1f-%.1f\n",
				current.StartTime, current.EndTime, next.StartTime, next.EndTime, current.StartTime, current.EndTime)
		} else {
			// No overlap, add current to merged and start new current
			merged = append(merged, current)
			current = next
		}
	}

	// Add the last range
	merged = append(merged, current)

	return merged
}

// findCharactersInTimeRange finds all characters that appear in a specific time range
func (s *characterService) findCharactersInTimeRange(videoID uuid.UUID, startTime, endTime float64) ([]characterModel.VideoSceneCharacter, error) {
	// Query to find all character appearances that overlap with the time range
	query := `
		SELECT DISTINCT
			ca.character_id,
			c.name as character_name,
			COALESCE(c.avatar, '') as character_avatar,
			AVG(ca.confidence) as confidence,
			ca.start_time,
			ca.end_time
		FROM character_appearances ca
		JOIN characters c ON c.id = ca.character_id AND c.is_active = true
		WHERE ca.video_id = $1
		AND ca.start_time <= $3
		AND ca.end_time >= $2
		GROUP BY ca.character_id, c.name, c.avatar, ca.start_time, ca.end_time
		ORDER BY ca.start_time
	`

	rows, err := s.sc.DB().Raw(query, videoID, startTime, endTime).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var characters []characterModel.VideoSceneCharacter

	for rows.Next() {
		var characterID uuid.UUID
		var characterName, characterAvatar string
		var confidence, charStartTime, charEndTime float64

		err := rows.Scan(&characterID, &characterName, &characterAvatar, &confidence, &charStartTime, &charEndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		character := characterModel.VideoSceneCharacter{
			CharacterID:     characterID,
			CharacterName:   characterName,
			CharacterAvatar: characterAvatar,
			Confidence:      confidence,
			StartTime:       charStartTime,
			EndTime:         charEndTime,
		}

		characters = append(characters, character)
		fmt.Printf("[DEBUG] Found character %s in range %.1f-%.1f (original: %.1f-%.1f)\n",
			characterName, startTime, endTime, charStartTime, charEndTime)
	}

	if len(characters) > 1 {
		intersectionStart := startTime
		intersectionEnd := endTime

		for _, char := range characters {
			if char.StartTime > intersectionStart {
				intersectionStart = char.StartTime
			}
			if char.EndTime < intersectionEnd {
				intersectionEnd = char.EndTime
			}
		}

		if intersectionStart <= intersectionEnd {
			for i := range characters {
				characters[i].StartTime = intersectionStart
				characters[i].EndTime = intersectionEnd
			}
			fmt.Printf("[DEBUG] Updated all characters to intersection time: %.1f-%.1f\n", intersectionStart, intersectionEnd)
		}
	}

	return characters, nil
}

// sceneContainsAllRequiredCharacters checks if a scene contains all required characters
func (s *characterService) sceneContainsAllRequiredCharacters(sceneCharacters []characterModel.VideoSceneCharacter, requiredCharacters []uuid.UUID) bool {
	if len(requiredCharacters) == 0 {
		return true
	}

	// Create a set of character IDs in the scene
	sceneCharacterIDs := make(map[uuid.UUID]bool)
	for _, char := range sceneCharacters {
		sceneCharacterIDs[char.CharacterID] = true
	}

	for _, requiredID := range requiredCharacters {
		if !sceneCharacterIDs[requiredID] {
			return false
		}
	}

	return true
}

func formatSecondsToTime(seconds float64) string {
	totalSeconds := int(seconds)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}
