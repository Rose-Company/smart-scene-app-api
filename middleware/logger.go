package middleware

import (
	"bytes"
	"smart-scene-app-api/common"
	"smart-scene-app-api/server"
	logger2 "smart-scene-app-api/services/logger_v2"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(sc server.ServerContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		requestId := c.Writer.Header().Get("X-Request-Id")
		c.Set("X-Request-Id", requestId)
		log := logger2.NewApiRequestLogger(c)

		c.Next()

		userId := "Unknown"

		val, ok := c.Get(common.UserId)
		if ok {
			userId = val.(string)
		}

		logData := logger2.LogDataObject{
			RequestId:          requestId,
			LogType:            "",
			ActorId:            userId,
			Path:               c.Request.URL.Path,
			Method:             c.Request.Method,
			Status:             c.Writer.Status(),
			ApplicationLatency: time.Since(now).Milliseconds(),
			RequestBody:        nil,
			Time:               now.Unix(),
		}

		if c.Writer.Status() != 200 {
			logData.Message = blw.body.String()
			log.Error(logData.ToString())
		} else {
			log.Info(logData.ToString())
		}
	}
}
