package tags

import (
	"errors"
	"net/http"
	"smart-scene-app-api/common"
	tagModels "smart-scene-app-api/internal/models/tag"
	tagService "smart-scene-app-api/internal/services/tag"
	"smart-scene-app-api/server"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	ctx        server.ServerContext
	tagService tagService.Service
}

func NewTagHandler(ctx server.ServerContext) *TagHandler {
	return &TagHandler{
		ctx:        ctx,
		tagService: tagService.NewTagService(ctx),
	}
}

func (h *TagHandler) GetTagsHierarchy(c *gin.Context) {
	var req tagModels.TagFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	req.VerifyPaging()

	response, err := h.tagService.GetTagsHierarchy(req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TagHandler) GetTags(c *gin.Context) {
	var req tagModels.TagFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	req.VerifyPaging()

	response, err := h.tagService.GetTagsList(req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TagHandler) GetTagsByPosition(c *gin.Context) {
	positionCode := c.Param("position_code")
	if positionCode == "" {
		common.AbortWithError(c, errors.New("position code is required"))
		return
	}

	var req tagModels.TagFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	req.PositionCode = positionCode
	req.VerifyPaging()

	response, err := h.tagService.GetTagsByPosition(req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TagHandler) GetTagsByCategory(c *gin.Context) {
	categoryCode := c.Param("category_code")
	if categoryCode == "" {
		common.AbortWithError(c, errors.New("category code is required"))
		return
	}

	var req tagModels.TagFilterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	req.CategoryCode = categoryCode
	req.VerifyPaging()

	response, err := h.tagService.GetTagsByCategory(req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TagHandler) GetTagUsageStats(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	sortBy := c.DefaultQuery("sort_by", "usage_count")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	response, err := h.tagService.GetTagUsageStats(limit, sortBy, sortOrder)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
