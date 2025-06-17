package tags

import (
	"errors"
	"net/http"
	"smart-scene-app-api/common"
	tagModels "smart-scene-app-api/internal/models/tag"
	tagService "smart-scene-app-api/internal/services/tag"
	"smart-scene-app-api/server"

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

	response, err := h.tagService.GetTagsByPosition(c.Request.Context(), req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
