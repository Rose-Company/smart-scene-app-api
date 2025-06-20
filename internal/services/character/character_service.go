package character

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models"
	characterModel "smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"
	characterRepo "smart-scene-app-api/internal/repositories/character"
	"smart-scene-app-api/server"

	"fmt"

	"math"
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

	filter := func(tx *gorm.DB) {
		tx.Where("video_id = ?", uuidID).
			Preload("Character", "is_active = ?", true).
			Order("start_time ASC")
	}

	appearances, err := s.appearanceRepo.List(s.sc.Ctx(), models.QueryParams{}, filter)
	if err != nil {
		return nil, err
	}

	if len(appearances) == 0 {
		return &characterModel.VideoSceneListResponse{
			BaseListResponse: models.BaseListResponse{
				Total:    0,
				Page:     queryParams.Page,
				PageSize: queryParams.PageSize,
			},
			Items: []characterModel.VideoScene{},
		}, nil
	}

	// Find time segments that match include/exclude criteria
	validSegments := s.findValidTimeSegments(appearances, queryParams.IncludeCharacters, queryParams.ExcludeCharacters)

	// Convert segments to scenes
	scenes := s.convertSegmentsToScenes(validSegments, uuidID)

	// Apply pagination
	total := len(scenes)
	offset := (queryParams.Page - 1) * queryParams.PageSize
	limit := queryParams.PageSize

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

	return response, nil
}

func (s *characterService) filterScenesByCharacters(scenes []characterModel.VideoScene, includeCharacters, excludeCharacters []uuid.UUID) []characterModel.VideoScene {
	return []characterModel.VideoScene{}
}

// TimeSegment represents a time interval with characters present
type TimeSegment struct {
	StartTime  float64
	EndTime    float64
	Characters map[uuid.UUID]*characterModel.CharacterAppearance
}

// findValidTimeSegments finds time segments where include characters are present and exclude characters are absent
func (s *characterService) findValidTimeSegments(appearances []*characterModel.CharacterAppearance, includeCharacters, excludeCharacters []uuid.UUID) []TimeSegment {
	if len(includeCharacters) == 0 {
		return []TimeSegment{}
	}

	// Step 1: Create time intervals for each include character
	includeIntervals := make(map[uuid.UUID][]TimeSegment)
	for _, charID := range includeCharacters {
		intervals := []TimeSegment{}
		for _, appearance := range appearances {
			if appearance.CharacterID == charID && appearance.Character != nil {
				intervals = append(intervals, TimeSegment{
					StartTime: appearance.StartTime,
					EndTime:   appearance.EndTime,
					Characters: map[uuid.UUID]*characterModel.CharacterAppearance{
						charID: appearance,
					},
				})
			}
		}
		includeIntervals[charID] = intervals
	}

	// Step 2: Find intersection of all include character intervals
	var intersectionSegments []TimeSegment
	if len(includeIntervals) > 0 {
		// Start with first character's intervals
		firstCharID := includeCharacters[0]
		intersectionSegments = includeIntervals[firstCharID]

		// Intersect with each subsequent character's intervals
		for i := 1; i < len(includeCharacters); i++ {
			charID := includeCharacters[i]
			intersectionSegments = s.intersectTimeSegments(intersectionSegments, includeIntervals[charID])
		}
	}

	// Step 3: Create exclude intervals if any
	excludeIntervals := []TimeSegment{}
	for _, charID := range excludeCharacters {
		for _, appearance := range appearances {
			if appearance.CharacterID == charID && appearance.Character != nil {
				excludeIntervals = append(excludeIntervals, TimeSegment{
					StartTime: appearance.StartTime,
					EndTime:   appearance.EndTime,
					Characters: map[uuid.UUID]*characterModel.CharacterAppearance{
						charID: appearance,
					},
				})
			}
		}
	}

	// Step 4: Remove exclude intervals from intersection
	validSegments := s.subtractTimeSegments(intersectionSegments, excludeIntervals)

	// Step 5: Merge adjacent segments and add all character info
	finalSegments := s.mergeAndEnrichSegments(validSegments, appearances)

	return finalSegments
}

// intersectTimeSegments finds the intersection of two sets of time segments
func (s *characterService) intersectTimeSegments(segments1, segments2 []TimeSegment) []TimeSegment {
	var result []TimeSegment

	for _, seg1 := range segments1 {
		for _, seg2 := range segments2 {
			// Find overlap
			startTime := math.Max(seg1.StartTime, seg2.StartTime)
			endTime := math.Min(seg1.EndTime, seg2.EndTime)

			if startTime < endTime { // Valid overlap
				// Merge character information
				characters := make(map[uuid.UUID]*characterModel.CharacterAppearance)
				for charID, appearance := range seg1.Characters {
					characters[charID] = appearance
				}
				for charID, appearance := range seg2.Characters {
					characters[charID] = appearance
				}

				result = append(result, TimeSegment{
					StartTime:  startTime,
					EndTime:    endTime,
					Characters: characters,
				})
			}
		}
	}

	return result
}

// subtractTimeSegments removes exclude segments from include segments
func (s *characterService) subtractTimeSegments(includeSegments, excludeSegments []TimeSegment) []TimeSegment {
	result := includeSegments

	for _, excludeSegment := range excludeSegments {
		var newResult []TimeSegment

		for _, includeSegment := range result {
			// Check if there's overlap
			if includeSegment.EndTime <= excludeSegment.StartTime || includeSegment.StartTime >= excludeSegment.EndTime {
				// No overlap, keep the segment
				newResult = append(newResult, includeSegment)
			} else {
				// There's overlap, split the segment

				// Part before exclude segment
				if includeSegment.StartTime < excludeSegment.StartTime {
					newResult = append(newResult, TimeSegment{
						StartTime:  includeSegment.StartTime,
						EndTime:    excludeSegment.StartTime,
						Characters: includeSegment.Characters,
					})
				}

				// Part after exclude segment
				if includeSegment.EndTime > excludeSegment.EndTime {
					newResult = append(newResult, TimeSegment{
						StartTime:  excludeSegment.EndTime,
						EndTime:    includeSegment.EndTime,
						Characters: includeSegment.Characters,
					})
				}
			}
		}

		result = newResult
	}

	return result
}

// mergeAndEnrichSegments merges adjacent segments and enriches with all character information
func (s *characterService) mergeAndEnrichSegments(segments []TimeSegment, allAppearances []*characterModel.CharacterAppearance) []TimeSegment {
	if len(segments) == 0 {
		return segments
	}

	// Sort segments by start time
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].StartTime < segments[j].StartTime
	})

	var enrichedSegments []TimeSegment

	for _, segment := range segments {
		// Find all characters present in this time segment
		characters := make(map[uuid.UUID]*characterModel.CharacterAppearance)

		for _, appearance := range allAppearances {
			if appearance.Character != nil &&
				appearance.StartTime < segment.EndTime &&
				appearance.EndTime > segment.StartTime {
				// Character appears in this segment
				characters[appearance.CharacterID] = appearance
			}
		}

		enrichedSegments = append(enrichedSegments, TimeSegment{
			StartTime:  segment.StartTime,
			EndTime:    segment.EndTime,
			Characters: characters,
		})
	}

	return enrichedSegments
}

// convertSegmentsToScenes converts time segments to video scenes
func (s *characterService) convertSegmentsToScenes(segments []TimeSegment, videoID uuid.UUID) []characterModel.VideoScene {
	var scenes []characterModel.VideoScene

	for i, segment := range segments {
		var sceneCharacters []characterModel.VideoSceneCharacter

		for _, appearance := range segment.Characters {
			sceneCharacters = append(sceneCharacters, characterModel.VideoSceneCharacter{
				CharacterID:     appearance.CharacterID,
				CharacterName:   appearance.Character.Name,
				CharacterAvatar: appearance.Character.Avatar,
				Confidence:      appearance.Confidence,
				StartTime:       math.Max(appearance.StartTime, segment.StartTime),
				EndTime:         math.Min(appearance.EndTime, segment.EndTime),
				StartFrame:      appearance.StartFrame,
				EndFrame:        appearance.EndFrame,
			})
		}

		scene := characterModel.VideoScene{
			VideoID:            videoID,
			SceneID:            fmt.Sprintf("segment_%d_%.1f_%.1f", i+1, segment.StartTime, segment.EndTime),
			StartTime:          segment.StartTime,
			EndTime:            segment.EndTime,
			Duration:           segment.EndTime - segment.StartTime,
			StartFrame:         0, // Will be calculated from start time if needed
			EndFrame:           0, // Will be calculated from end time if needed
			CharacterCount:     len(sceneCharacters),
			Characters:         sceneCharacters,
			StartTimeFormatted: formatSecondsToTime(segment.StartTime),
			EndTimeFormatted:   formatSecondsToTime(segment.EndTime),
		}

		scenes = append(scenes, scene)
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
