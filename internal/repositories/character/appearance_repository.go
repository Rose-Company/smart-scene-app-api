package character

import (
	"context"
	"encoding/json"
	"fmt"
	"smart-scene-app-api/internal/models/character"
	"smart-scene-app-api/internal/repositories"

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

	// Simple query that groups by time segments and filters characters
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
		WHERE ca.video_id = $1
		GROUP BY start_time, end_time, duration
		HAVING 
			-- Có ít nhất 1 nhân vật trong include list
			SUM(CASE WHEN ca.character_id IN (%s) THEN 1 ELSE 0 END) > 0
			%s
		ORDER BY start_time
	`

	// Build include characters condition
	includePlaceholders := make([]string, len(includeCharacters))
	args := []interface{}{videoID} // First arg for video_id

	for i, id := range includeCharacters {
		includePlaceholders[i] = fmt.Sprintf("$%d", len(args)+1)
		args = append(args, id)
	}
	includeCondition := strings.Join(includePlaceholders, ",")

	// Build exclude characters condition if any
	excludeCondition := ""
	if len(excludeCharacters) > 0 {
		excludePlaceholders := make([]string, len(excludeCharacters))
		for i, id := range excludeCharacters {
			excludePlaceholders[i] = fmt.Sprintf("$%d", len(args)+1)
			args = append(args, id)
		}
		excludeCondition = fmt.Sprintf("AND SUM(CASE WHEN ca.character_id IN (%s) THEN 1 ELSE 0 END) = 0", strings.Join(excludePlaceholders, ","))
	}

	// Format the final query
	finalQuery := fmt.Sprintf(query, includeCondition, excludeCondition)

	fmt.Printf("[DEBUG] Repository SQL Query Parameters - VideoID: %s, Include: %v, Exclude: %v\n",
		videoID, includeCharacters, excludeCharacters)
	fmt.Printf("[DEBUG] Final Query: %s\n", finalQuery)
	fmt.Printf("[DEBUG] Args: %v\n", args)

	// Execute query
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
			fmt.Printf("[DEBUG] Repository error scanning row: %v\n", err)
			return nil, fmt.Errorf("failed to scan query result row: %w", err)
		}

		fmt.Printf("[DEBUG] Raw JSON from database: %s\n", charactersJSON)

		// Parse characters JSON
		var charactersData []map[string]interface{}
		if err := json.Unmarshal([]byte(charactersJSON), &charactersData); err != nil {
			fmt.Printf("[DEBUG] Repository error parsing characters JSON: %v\n", err)
			fmt.Printf("[DEBUG] Problematic JSON string: %s\n", charactersJSON)
			return nil, fmt.Errorf("failed to parse characters JSON: %w", err)
		}

		// Convert to characters in segment
		var characters []character.CharacterInSegment
		for _, char := range charactersData {
			// Check if character data is valid
			if char == nil {
				fmt.Printf("[DEBUG] Skipping null character entry\n")
				continue
			}

			characterIDStr, ok := char["character_id"].(string)
			if !ok {
				fmt.Printf("[DEBUG] character_id is not a string: %v\n", char["character_id"])
				return nil, fmt.Errorf("failed to parse character_id as string")
			}

			charID, err := uuid.Parse(characterIDStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse character ID: %w", err)
			}

			confidence, ok := char["confidence"].(float64)
			if !ok {
				fmt.Printf("[DEBUG] confidence is not a float64: %v\n", char["confidence"])
				return nil, fmt.Errorf("failed to parse confidence value")
			}

			characterName, ok := char["character_name"].(string)
			if !ok {
				fmt.Printf("[DEBUG] character_name is not a string: %v\n", char["character_name"])
				return nil, fmt.Errorf("failed to parse character name")
			}

			characterAvatar, ok := char["character_avatar"].(string)
			if !ok {
				// Avatar might be null, handle gracefully
				characterAvatar = ""
				fmt.Printf("[DEBUG] character_avatar is null or not a string, using empty string\n")
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
		fmt.Printf("[DEBUG] Repository created time segment: %.1f-%.1f (%.1fs) with %d characters\n",
			startTime, endTime, duration, len(characters))
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return results, nil
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
