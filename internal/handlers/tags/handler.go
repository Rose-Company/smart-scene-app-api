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

// GetTagsByPosition godoc
// @Summary      Get tags by position
// @Description  Get tags by position
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        position_code  path      string  true  "Position code"
// @Param        query  query      tagModels.TagFilterRequest  true  "Query parameters"
// @Success      200  {object}  common.Response{data=tagModels.TagListResponse}  "Tags retrieved successfully"
// @Failure      400  {object}  common.Response  "Bad request"
// @Failure      401  {object}  common.Response  "Unauthorized"
// @Failure      500  {object}  common.Response  "Internal server error"
// @Router       /api/v1/tags/position/{position_code} [get]
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
