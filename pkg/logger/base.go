package l

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EncoderName = "flat-encoder"
)

var logEnv string

// ShortColorCallerEncoder encodes caller information with sort path filename and enable color.
func ShortColorCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	callerStr := caller.TrimmedPath() + ":" + strconv.Itoa(caller.Line)
	enc.AppendString(callerStr)
}

// DefaultConsoleEncoderConfig ...
var DefaultConsoleEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   ShortColorCallerEncoder,
}

func init() {
	logEnv = os.Getenv("ENVIRONMENT")
	var logLevel = os.Getenv("LOG_LEVEL")
	var err error

	var logConfig zap.Config
	if strings.ToUpper(logEnv) == "STAG" {
		logConfig = zap.Config{
			Level:            zap.NewAtomicLevel(),
			Development:      false,
			Encoding:         EncoderName,
			EncoderConfig:    DefaultConsoleEncoderConfig,
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		logConfig = zap.NewProductionConfig()
		logConfig.Level.SetLevel(zapcore.WarnLevel)
	}

	if len(logLevel) > 0 {
		if err := logConfig.Level.UnmarshalText([]byte(logLevel)); err != nil {
			log.Fatal("Cannot setting log level", logLevel, err)
		}
	}

	l, err = logConfig.Build()
	if err != nil {
		log.Fatal("Cannot init logger", err)
	}

}

var l *zap.Logger

func New() *zap.Logger {
	return l.Named("")
}

func NewWithName(name string) *zap.Logger {
	return l.Named(name)
}

// Short-hand functions for logging.
var (
	Any        = zap.Any
	Bool       = zap.Bool
	Duration   = zap.Duration
	Float64    = zap.Float64
	Int        = zap.Int
	Int64      = zap.Int64
	Skip       = zap.Skip
	String     = zap.String
	Stringer   = zap.Stringer
	Time       = zap.Time
	Uint       = zap.Uint
	Uint32     = zap.Uint32
	Uint64     = zap.Uint64
	Uintptr    = zap.Uintptr
	ByteString = zap.ByteString
)

// Object ...
func Object(key string, val interface{}) zapcore.Field {
	//return zap.Any(key, val)
	return zap.Stringer(key, Dump(val))
}

type dd struct {
	v interface{}
}

func (d dd) String() string {
	return pp.Sprint(d.v)
}

// Dump renders object for debugging
func Dump(v interface{}) fmt.Stringer {
	return dd{v}
}

// Error wraps error for zap.Error.
func Error(err error) zapcore.Field {
	if err == nil {
		return Skip()
	}
	return String("error", err.Error())
}
