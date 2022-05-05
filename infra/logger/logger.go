package logger

import (
	"context"
	"io"

	"github.com/diegoclair/go-boilerplate/util/config"
)

// Level represents a logging level
type Level uint8

// Logging level standard values
const (
	PANIC Level = iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
)

// Logger is the default application logger definition
type Logger interface {
	Level() Level
	SetLevel(level Level)

	AppName() string
	SetAppName(name string)

	RecordID() string
	SetRecordID(id string)

	Output() io.Writer
	SetOutput(output io.Writer)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	InfoWriter() io.Writer
	ErrorWriter() io.Writer
	FatalWriter() io.Writer

	WithFields(map[string]interface{}) Logger
	Fields() map[string]interface{}

	NewSessionLogger(appContext context.Context) (context.Context, Logger)
}

// New returns a new logger instance
func New(cfg config.LogConfig, appName string) Logger {
	logger := newLogrusLogger(cfg)

	logger.SetAppName(appName)

	if cfg.Debug {
		logger.SetLevel(DEBUG)
	} else {
		logger.SetLevel(INFO)
	}

	return logger
}
