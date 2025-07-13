package character

import (
	"context"
	"encoding/json"
	"fmt"
	"smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"

	"math"
	"sort"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppearanceRepository interface {
	repositories.BaseRepository[character.CharacterAppearance]
	FindTimeSegmentsWithCharacters(ctx context.Context, videoID uuid.UUID, includeCharacters, excludeCharacters []uuid.UUID) ([]character.TimeSegmentResult, error)
}

type appearanceRepository struct {
	repositories.BaseRepository[character.CharacterAppearance]
	db *gorm.DB
}

func NewAppearanceRepository(db *gorm.DB) AppearanceRepository {
	return &appearanceRepository{
		BaseRepository: repositories.NewBaseRepository[character.CharacterAppearance](db),
		db:             db,
	}
}

func (r *appearanceRepository) FindTimeSegmentsWithCharacters(ctx context.Context, videoID uuid.UUID, includeCharacters, excludeCharacters []uuid.UUID) ([]character.TimeSegmentResult, error) {
	if len(includeCharacters) == 0 {
		return []character.TimeSegmentResult{}, nil
	}

	fmt.Printf("[DEBUG] Repository SQL Query Parameters - VideoID: %s, Include: %v, Exclude: %v\n",
		videoID, includeCharacters, excludeCharacters)

	// Step 1: Get all segments with include characters
	includeSegments, err := r.getSegmentsWithCharacters(ctx, videoID, includeCharacters)
	if err != nil {
		return nil, fmt.Errorf("failed to get include segments: %w", err)
	}

	// Step 2: If no exclude characters, return include segments as-is
	if len(excludeCharacters) == 0 {
		return includeSegments, nil
	}

	// Step 3: Get all segments with exclude characters
	excludeSegments, err := r.getSegmentsWithCharacters(ctx, videoID, excludeCharacters)
	if err != nil {
		return nil, fmt.Errorf("failed to get exclude segments: %w", err)
	}

	// Step 4: Cut overlap between include and exclude segments
	filteredSegments := r.cutOverlapFromSegments(includeSegments, excludeSegments)

	fmt.Printf("[DEBUG] Final filtered segments: %d\n", len(filteredSegments))
	return filteredSegments, nil
}

// getSegmentsWithCharacters gets time segments containing specific characters
func (r *appearanceRepository) getSegmentsWithCharacters(ctx context.Context, videoID uuid.UUID, characterIDs []uuid.UUID) ([]character.TimeSegmentResult, error) {
	if len(characterIDs) == 0 {
		return []character.TimeSegmentResult{}, nil
	}

	query := `
		SELECT 
			start_time,
			end_time,
			duration,
			COUNT(*) as character_count,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'character_id', ca.character_id::text,
					'character_name', COALESCE(c.name, ''),
					'character_avatar', COALESCE(c.avatar, ''),
					'confidence', COALESCE(ca.confidence, 0.0)
				)
			) as characters
		FROM character_appearances ca
		JOIN characters c ON c.id = ca.character_id AND c.is_active = true
		WHERE ca.video_id = $1 AND ca.character_id IN (%s)
		GROUP BY start_time, end_time, duration
		ORDER BY start_time
	`

	// Build placeholders for character IDs
	placeholders := make([]string, len(characterIDs))
	args := []interface{}{videoID}
	for i, id := range characterIDs {
		placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
		args = append(args, id)
	}

	finalQuery := fmt.Sprintf(query, strings.Join(placeholders, ","))
	fmt.Printf("[DEBUG] Query: %s\n", finalQuery)
	fmt.Printf("[DEBUG] Args: %v\n", args)

	rows, err := r.db.Raw(finalQuery, args...).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []character.TimeSegmentResult

	for rows.Next() {
		var startTime, endTime, duration float64
		var characterCount int
		var charactersJSON string

		err := rows.Scan(&startTime, &endTime, &duration, &characterCount, &charactersJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan query result row: %w", err)
		}

		// Parse characters JSON
		var charactersData []map[string]interface{}
		if err := json.Unmarshal([]byte(charactersJSON), &charactersData); err != nil {
			return nil, fmt.Errorf("failed to parse characters JSON: %w", err)
		}

		// Convert to characters in segment
		var characters []character.CharacterInSegment
		for _, char := range charactersData {
			if char == nil {
				continue
			}

			characterIDStr, ok := char["character_id"].(string)
			if !ok {
				return nil, fmt.Errorf("failed to parse character_id as string")
			}

			charID, err := uuid.Parse(characterIDStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse character ID: %w", err)
			}

			confidence, ok := char["confidence"].(float64)
			if !ok {
				return nil, fmt.Errorf("failed to parse confidence value")
			}

			characterName, ok := char["character_name"].(string)
			if !ok {
				return nil, fmt.Errorf("failed to parse character name")
			}

			characterAvatar, ok := char["character_avatar"].(string)
			if !ok {
				characterAvatar = ""
			}

			characters = append(characters, character.CharacterInSegment{
				CharacterID:     charID,
				CharacterName:   characterName,
				CharacterAvatar: characterAvatar,
				Confidence:      confidence,
			})
		}

		result := character.TimeSegmentResult{
			StartTime:       startTime,
			EndTime:         endTime,
			Duration:        duration,
			TotalCharacters: characterCount,
			Characters:      characters,
		}

		results = append(results, result)
	}

	return results, nil
}

// cutOverlapFromSegments cuts overlapping portions from include segments where exclude segments exist
func (r *appearanceRepository) cutOverlapFromSegments(includeSegments, excludeSegments []character.TimeSegmentResult) []character.TimeSegmentResult {
	var filteredSegments []character.TimeSegmentResult

	for _, includeSegment := range includeSegments {
		fmt.Printf("[DEBUG] Processing include segment: %.1f-%.1f\n", includeSegment.StartTime, includeSegment.EndTime)

		// Find all exclude segments that overlap with this include segment
		var overlappingExcludes []character.TimeSegmentResult
		for _, excludeSegment := range excludeSegments {
			if r.segmentsOverlap(includeSegment, excludeSegment) {
				overlappingExcludes = append(overlappingExcludes, excludeSegment)
				fmt.Printf("[DEBUG] Found overlapping exclude segment: %.1f-%.1f\n", excludeSegment.StartTime, excludeSegment.EndTime)
			}
		}

		// If no overlapping excludes, keep the include segment as-is
		if len(overlappingExcludes) == 0 {
			filteredSegments = append(filteredSegments, includeSegment)
			fmt.Printf("[DEBUG] No overlaps, keeping segment: %.1f-%.1f\n", includeSegment.StartTime, includeSegment.EndTime)
			continue
		}

		// Cut the include segment based on overlapping excludes
		cutSegments := r.cutSegmentWithExcludes(includeSegment, overlappingExcludes)
		filteredSegments = append(filteredSegments, cutSegments...)

		for _, cutSegment := range cutSegments {
			fmt.Printf("[DEBUG] Cut segment result: %.1f-%.1f\n", cutSegment.StartTime, cutSegment.EndTime)
		}
	}

	return filteredSegments
}

// segmentsOverlap checks if two time segments overlap
func (r *appearanceRepository) segmentsOverlap(seg1, seg2 character.TimeSegmentResult) bool {
	return seg1.StartTime < seg2.EndTime && seg2.StartTime < seg1.EndTime
}

// cutSegmentWithExcludes cuts an include segment by removing overlapping portions with exclude segments
func (r *appearanceRepository) cutSegmentWithExcludes(includeSegment character.TimeSegmentResult, excludeSegments []character.TimeSegmentResult) []character.TimeSegmentResult {
	var results []character.TimeSegmentResult
	currentStart := includeSegment.StartTime
	currentEnd := includeSegment.EndTime

	// Sort exclude segments by start time
	sort.Slice(excludeSegments, func(i, j int) bool {
		return excludeSegments[i].StartTime < excludeSegments[j].StartTime
	})

	for _, exclude := range excludeSegments {
		// Calculate overlap
		overlapStart := math.Max(currentStart, exclude.StartTime)
		overlapEnd := math.Min(currentEnd, exclude.EndTime)

		// If there's actual overlap
		if overlapStart < overlapEnd {
			// Add segment before overlap (if any)
			if currentStart < overlapStart {
				beforeSegment := character.TimeSegmentResult{
					StartTime:       currentStart,
					EndTime:         overlapStart,
					Duration:        overlapStart - currentStart,
					TotalCharacters: includeSegment.TotalCharacters,
					Characters:      includeSegment.Characters,
				}
				results = append(results, beforeSegment)
				fmt.Printf("[DEBUG] Added before-overlap segment: %.1f-%.1f\n", beforeSegment.StartTime, beforeSegment.EndTime)
			}

			// Move current start to after the overlap
			currentStart = overlapEnd
		}
	}

	// Add remaining segment after all overlaps (if any)
	if currentStart < currentEnd {
		afterSegment := character.TimeSegmentResult{
			StartTime:       currentStart,
			EndTime:         currentEnd,
			Duration:        currentEnd - currentStart,
			TotalCharacters: includeSegment.TotalCharacters,
			Characters:      includeSegment.Characters,
		}
		results = append(results, afterSegment)
		fmt.Printf("[DEBUG] Added after-overlap segment: %.1f-%.1f\n", afterSegment.StartTime, afterSegment.EndTime)
	}

	return results
}

// joinStrings joins string slices with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
