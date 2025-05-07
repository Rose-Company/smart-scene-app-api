package handlers

import (
	"smart-scene-app-api/common"
	"smart-scene-app-api/server"

	"github.com/gin-gonic/gin"
)

func Check(sc server.ServerContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(common.SUCCESS_STATUS, common.ResponseSuccess(common.REQUEST_SUCCESS, "", "success"))
	}
}
