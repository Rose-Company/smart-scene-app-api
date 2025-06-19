package character

import (
	"net/http"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models/character"

	"github.com/gin-gonic/gin"
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
	videoID := c.Param("video_id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'video_id' parameter is missing or empty",
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
