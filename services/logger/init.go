package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

type Logger interface {
	Run() error
	Stop() <-chan bool
	GetPrefix() string
	Get() interface{}
}

type Loggers interface {
	Info() *log.Logger
	Warning() *log.Logger
	Error() *log.Logger
	ConfigureRequestId(requestId string)
}

type loggers struct {
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

func (l *loggers) Info() *log.Logger {
	return l.info
}

func (l *loggers) Warning() *log.Logger {
	return l.warning
}

func (l *loggers) Error() *log.Logger {
	return l.error
}

type logger struct {
	prefix string
	logs   *loggers
}

func NewLogger(prefix string) *logger {
	return &logger{prefix: prefix}
}

func (l *logger) GetPrefix() string {
	return l.prefix
}

func (l *logger) Run() error {
	if err := l.configure(); err != nil {
		return err
	}

	return nil
}

func (l *logger) configure() error {
	info := log.New(os.Stdout, Green+"[INFO]: "+Reset, log.Ldate|log.Ltime|log.Lshortfile)
	warning := log.New(os.Stdout, Yellow+"[WARN]: "+Reset, log.Ldate|log.Ltime|log.Lshortfile)
	error := log.New(os.Stdout, Red+"[ERROR]: "+Reset, log.Ldate|log.Ltime|log.Lshortfile)

	logs := &loggers{
		info:    info,
		warning: warning,
		error:   error,
	}

	l.logs = logs

	return nil
}

func (l *logger) Stop() <-chan bool {
	stop := make(chan bool, 1)
	go func() {
		stop <- true
	}()
	l.logs.info.Println("Logger is stopped")
	return stop
}

func (l *logger) Get() *loggers {
	return l.logs
}

func (l *loggers) ConfigureRequestId(requestId string) {
	l.info = log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.warning = log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[WARN]: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.error = log.New(os.Stdout, fmt.Sprintf("[REQUEST-ID-%v] ", requestId)+"[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}
