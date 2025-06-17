package video

import (
	"net/http"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models/video"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllVideos godoc
// @Summary      Get all videos
// @Description  Retrieve a list of all videos
// @Tags         videos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.Response{data=[]video.Video}  "List of videos"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos [get]
func (h *Handler) GetAllVideos(c *gin.Context) {
	var queryParams video.VideoFilterAndPagination
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid query parameters",
			ErrorDetail: err.Error(),
		})
		return
	}

	videos, err := h.service.Video.GetAllVideos(queryParams)
	if err != nil {
		h.logger.Error("Failed to get videos: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve videos",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Videos retrieved successfully",
		Data:    videos,
	})
}

// GetVideoByID godoc
// @Summary      Get video by ID
// @Description  Retrieve a video by its ID
// @Tags         videos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Video ID"
// @Success      200  {object}  common.Response{data=video.Video}  "Video details"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      404  {object}  common.Response  "Video not found"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos/{id} [get]
func (h *Handler) GetVideoByID(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	video, err := h.service.Video.GetVideoByID(videoID)
	if err != nil {
		if err == common.ErrVideoNotFound {
			c.JSON(http.StatusNotFound, common.Response{
				Message:     "Video not found",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to get video by ID: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve video",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Video retrieved successfully",
		Data:    video,
	})
}

// CreateVideo godoc
// @Summary      Create a new video
// @Description  Create a new video with the provided details
// @Tags         videos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        video  body      video.Video  true  "Video details"
// @Success      201  {object}  common.Response{data=video.Video}  "Video created successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos [post]
func (h *Handler) CreateVideo(c *gin.Context) {
	var newVideo video.Video
	if err := c.ShouldBindJSON(&newVideo); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid video data",
			ErrorDetail: err.Error(),
		})
		return
	}

	createdVideo, err := h.service.Video.CreateVideo(newVideo)
	if err != nil {
		h.logger.Error("Failed to create video: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to create video",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.Response{
		Message: "Video created successfully",
		Data:    createdVideo,
	})
}

// UpdateVideo godoc
// @Summary      Update an existing video
// @Description  Update a video by its ID with the provided details
// @Tags         videos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Video ID"
// @Param        video  body      video.Video  true  "Updated video details"
// @Success      200  {object}  common.Response{data=video.Video}  "Video updated successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      404  {object}  common.Response  "Video not found"
// @Failure      500  {object}  common.Response  "Internal server error"
func (h *Handler) UpdateVideo(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	var updatedVideo video.Video
	if err := c.ShouldBindJSON(&updatedVideo); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid video data",
			ErrorDetail: err.Error(),
		})
		return
	}

	video, err := h.service.Video.UpdateVideo(videoID, updatedVideo)
	if err != nil {
		if err == common.ErrVideoNotFound {
			c.JSON(http.StatusNotFound, common.Response{
				Message:     "Video not found",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to update video: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to update video",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Video updated successfully",
		Data:    video,
	})
}

// DeleteVideo godoc
// @Summary      Delete a video
// @Description  Delete a video by its ID
// @Tags         videos
// @Security     BearerAuth
// @Param        id  path      string  true  "Video ID"
// @Success      204  {object}  common.Response  "Video deleted successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      404  {object}  common.Response  "Video not found"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos/{id} [delete]
func (h *Handler) DeleteVideo(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	if err := h.service.Video.DeleteVideo(videoID); err != nil {
		if err == common.ErrVideoNotFound {
			c.JSON(http.StatusNotFound, common.Response{
				Message:     "Video not found",
				ErrorDetail: err.Error(),
			})
			return
		}
		h.logger.Error("Failed to delete video: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to delete video",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, common.Response{
		Message: "Video deleted successfully",
	})
}

// GetVideosListing godoc
// @Summary      Get videos listing with search and filters
// @Description  Retrieve a paginated list of videos with search by title/character and tag filters
// @Tags         videos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page          query     int      false  "Page number (default: 1)"
// @Param        page_size     query     int      false  "Page size (default: 10)"
// @Param        sort          query     string   false  "Sort field and order (e.g., created_at.desc)"
// @Param        search        query     string   false  "Search by video title"
// @Param        character_names query   []string false  "Search by character names"
// @Param        character_ids query     []string false  "Search by character IDs"
// @Param        status        query     string   false  "Filter by status (pending, processing, completed, failed)"
// @Param        tag_ids       query     []int    false  "Filter by tag IDs"
// @Param        tag_codes     query     []string false  "Filter by tag codes (male, female, etc.)"
// @Param        date_from     query     string   false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to       query     string   false  "Filter to date (YYYY-MM-DD)"
// @Param        min_duration  query     int      false  "Minimum duration in seconds"
// @Param        max_duration  query     int      false  "Maximum duration in seconds"
// @Param        has_characters query    bool     false  "Filter videos with/without characters"
// @Param        created_by    query     string   false  "Filter by creator UUID"
// @Success      200  {object}  common.Response{data=video.VideoListResponse}  "Videos listing with filters"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos/listing [get]
func (h *Handler) GetVideosListing(c *gin.Context) {
	var req video.VideoListingRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Invalid query parameters",
			ErrorDetail: err.Error(),
		})
		return
	}

	// Verify paging parameters
	req.VerifyPaging()

	// Get videos listing from service with all filters
	response, err := h.service.Video.GetVideosListing(req)
	if err != nil {
		h.logger.Error("Failed to get videos listing: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve videos listing",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Videos listing retrieved successfully",
		Data:    response,
	})
}

// GetVideoSearchSuggestions godoc
// @Summary      Get search suggestions for videos and characters
// @Description  Get autocomplete suggestions for video titles and character names
// @Tags         videos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        q        query    string  true   "Search query"
// @Param        limit    query    int     false  "Maximum suggestions to return (default: 10)"
// @Success      200  {object}  common.Response{data=video.VideoSearchSuggestionResponse}  "Search suggestions"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/videos/search/suggestions [get]
func (h *Handler) GetVideoSearchSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Search query is required",
			ErrorDetail: "The 'q' parameter is missing or empty",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// Get search suggestions from service
	response, err := h.service.Video.GetVideoSearchSuggestions(query, limit)
	if err != nil {
		h.logger.Error("Failed to get search suggestions: " + err.Error())
		c.JSON(http.StatusInternalServerError, common.Response{
			Message:     "Failed to retrieve search suggestions",
			ErrorDetail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Message: "Search suggestions retrieved successfully",
		Data:    response,
	})
}
