package middleware

import (
	"fmt"
	"runtime/debug"
	"smart-scene-app-api/common"
	"smart-scene-app-api/server"

	"github.com/gin-gonic/gin"
)

func Recovery(sc server.ServerContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			c.Header("Content-Type", "application/json")

			logger := sc.GetLogger()

			if err := recover(); err != nil {
				logger.Error().Println(err)
				logger.Error().Println(string(debug.Stack()))
				c.AbortWithStatusJSON(common.SERVER_ERROR_STATUS, common.BaseResponse(common.REQUEST_FAILED, "Đã xảy ra lỗi, xin vui lòng thử lại", fmt.Sprintf("%v", err), nil))
			}
		}()

		c.Next()
	}
}
