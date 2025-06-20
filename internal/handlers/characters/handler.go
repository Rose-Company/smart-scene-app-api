package character

import (
	"net/http"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models/character"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCharactersByVideoID godoc
// @Summary      Get characters by video ID
// @Description  Retrieve a list of characters that appear in a specific video with appearance statistics
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        video_id  path      string  true  "Video ID"
// @Param        page      query     int     false "Page number (default: 1)"
// @Param        page_size query     int     false "Page size (default: 10, max: 100)"
// @Param        character_name query string false "Filter by character name"
// @Param        min_confidence query number false "Minimum confidence threshold"
// @Param        min_appearances query int   false "Minimum number of appearances"
// @Param        sort      query     string  false "Sort by: appearance_count.desc, total_duration.desc, first_appearance.asc, character_name.asc, confidence.desc"
// @Success      200  {object}  common.Response{data=character.VideoCharacterListResponse}  "Characters retrieved successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos/{video_id}/characters [get]
func (h *Handler) GetCharactersByVideoID(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	var queryParams character.VideoCharacterFilterAndPagination
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid query parameters",
			ErrorDetail: err.Error(),
		})
		return
	}

	characters, err := h.service.Character.GetCharactersByVideoID(videoID, queryParams)
	if err != nil {
		if err == common.ErrInvalidUUID {
			c.JSON(http.StatusBadRequest, common.Response{
				Message:     "Invalid video ID format",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get characters by video ID: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve characters",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Characters retrieved successfully",
		Data:    characters,
	})
}

// GetVideoScenesWithCharacters godoc
// @Summary      Get video scenes with character filtering
// @Description  Retrieve scenes from a video with include/exclude character filtering
// @Tags         characters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        video_id  path      string  true  "Video ID"
// @Param        page      query     int     false "Page number (default: 1)"
// @Param        page_size query     int     false "Page size (default: 10, max: 100)"
// @Param        include_characters query []string false "Character IDs that MUST be present in scene"
// @Param        exclude_characters query []string false "Character IDs that must NOT be present in scene"
// @Param        min_duration query number false "Minimum scene duration in seconds"
// @Param        max_duration query number false "Maximum scene duration in seconds"
// @Param        min_confidence query number false "Minimum character confidence threshold"
// @Param        overlap_threshold query number false "Time overlap threshold for grouping scenes (default: 1.0 seconds)"
// @Success      200  {object}  common.Response{data=character.VideoSceneListResponse}  "Scenes retrieved successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos/{video_id}/scenes [get]
func (h *Handler) GetVideoScenesWithCharacters(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	var queryParams character.VideoSceneFilterAndPagination
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid query parameters",
			ErrorDetail: err.Error(),
		})
		return
	}

	// Parse include characters
	for _, charStr := range queryParams.IncludeCharactersStr {
		if charStr != "" {
			charUUID, err := uuid.Parse(strings.TrimSpace(charStr))
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Response{
					Message:     "Invalid include character UUID",
					ErrorDetail: "Character ID '" + charStr + "' is not a valid UUID",
				})
				return
			}
			queryParams.IncludeCharacters = append(queryParams.IncludeCharacters, charUUID)
		}
	}

	// Parse exclude characters
	for _, charStr := range queryParams.ExcludeCharactersStr {
		if charStr != "" {
			charUUID, err := uuid.Parse(strings.TrimSpace(charStr))
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Response{
					Message:     "Invalid exclude character UUID",
					ErrorDetail: "Character ID '" + charStr + "' is not a valid UUID",
				})
				return
			}
			queryParams.ExcludeCharacters = append(queryParams.ExcludeCharacters, charUUID)
		}
	}

	scenes, err := h.service.Character.GetVideoScenesWithCharacters(videoID, queryParams)
	if err != nil {
		if err == common.ErrInvalidUUID {
			c.JSON(http.StatusBadRequest, common.Response{
				Message:     "Invalid video ID format",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get video scenes: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve video scenes",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Video scenes retrieved successfully",
		Data:    scenes,
	})
}
