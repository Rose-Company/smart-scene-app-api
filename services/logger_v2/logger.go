package logger

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type loggers struct {
	info      *log.Logger
	warning   *log.Logger
	error     *log.Logger
	requestId string
}

func (l *loggers) Info(v ...any) {
	//lData := LogDataObject{
	//	RequestId: l.requestId,
	//	Message:   v,
	//	LogType:   "INFO",
	//}
	l.info.Println(v)
}

func (l *loggers) Warning(v ...any) {
	//lData := LogDataObject{
	//	RequestId: l.requestId,
	//	Message:   v,
	//	LogType:   "WARNING",
	//}
	l.info.Println(v)
}

func (l *loggers) Error(v ...any) {
	//lData := LogDataObject{
	//	RequestId: l.requestId,
	//	Message:   v,
	//	LogType:   "ERROR",
	//}
	l.info.Println("Error", v)
}

func NewApiRequestLogger(c *gin.Context) *loggers {
	//logStn := log.New(os.Stdout, "", 0)
	requestId := c.Writer.Header().Get("X-Request-Id")

	info := log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	warning := log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[WARN]: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorL := log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)

	logs := &loggers{
		info:      info,
		warning:   warning,
		error:     errorL,
		requestId: requestId,
	}
	return logs
}

func NewJobLogger(requestId string) *loggers {

	info := log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	warning := log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[WARN]: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorL := log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)

	logs := &loggers{
		info:      info,
		warning:   warning,
		error:     errorL,
		requestId: requestId,
	}
	return logs
}

type ApiLoggerInterface interface {
	Info(v ...any)
	Warning(v ...any)
	Error(v ...any)
}

type LogDataObject struct {
	RequestId          string `json:"request_id,omitempty"`
	LogType            string `json:"log_type,omitempty"`
	ActorId            string `json:"actor_id,omitempty"`
	Path               string `json:"path,omitempty"`
	Method             string `json:"method,omitempty"`
	Status             int    `json:"status,omitempty"`
	ApplicationLatency int64  `json:"application_latency,omitempty"`
	RequestBody        any    `json:"request_body,omitempty"`
	Message            any    `json:"message,omitempty"`
	Time               int64  `json:"time,omitempty"`
}

func (l LogDataObject) ToString() string {
	str, _ := json.Marshal(l)
	return string(str)
}
