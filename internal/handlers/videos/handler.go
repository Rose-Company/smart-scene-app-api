package video

import (
	"net/http"
	"smart-scene-app-api/common"
	"smart-scene-app-api/internal/models/video"

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

// GetVideoDetail godoc
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
func (h *Handler) GetVideoDetail(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, common.Response{
			Message:     "Video ID is required",
			ErrorDetail: "The 'id' parameter is missing or empty",
		})
		return
	}

	video, err := h.service.Video.GetVideoDetail(videoID)
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
