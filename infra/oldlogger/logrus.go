package oldlogger

import (
	"context"
	"fmt"
	"os"

	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/sirupsen/logrus"
)

// entryKey standard values
const (
	entryKeyAppName  = "app"
	entryKeyRecordID = "id"
)

func newLogrusLogger(cfg config.Config) Logger {
	logger := &LogrusLogger{cfg: cfg}

	if cfg.Log.LogToFile {
		file, err := os.Create(cfg.Log.Path)
		if err != nil {
			fmt.Printf("Error to create log file for library: %s\n", err.Error())
			panic(err)
		}
		logrus.SetOutput(file)
	}

	logrus.SetFormatter(&customJSONFormatter{cfg: cfg.Log})

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Errorf("Error obtaining host name: %v", err)
	}

	logger.Entry = logrus.WithFields(logrus.Fields{
		"hostname": hostname,
		"file":     "",
	})

	logger.SetAppName(cfg.App.Name)

	if cfg.Log.Debug {
		logger.SetLevel(DEBUG)
	} else {
		logger.SetLevel(INFO)
	}

	return logger
}

// LogrusLogger is the default application logger
type LogrusLogger struct {
	cfg config.Config
	*logrus.Entry
}

// NewSessionLogger returns an instance of log with session code field
func (l *LogrusLogger) NewSessionLogger(ctx context.Context) (context.Context, Logger) {

	sessionCode := ctx.Value(auth.SessionKey)
	if sessionCode == nil {
		return ctx, newLogrusLogger(l.cfg)
	}

	var sessionCodeKey = "logger-" + sessionCode.(string)
	vl := ctx.Value(sessionCodeKey)
	if vl != nil {
		return ctx, vl.(Logger)
	}

	instance := newLogrusLogger(l.cfg)
	instance = instance.WithFields(map[string]interface{}{
		"session_code": sessionCode.(string),
	})
	ctx = context.WithValue(ctx, sessionCodeKey, instance)
	return ctx, instance
}

func (l *LogrusLogger) AppName() string {

	appName, ok := l.Entry.Data[entryKeyAppName]
	if !ok {
		return ""
	}

	return appName.(string)
}

func (l *LogrusLogger) SetAppName(name string) {
	if name == "" {
		delete(l.Entry.Data, entryKeyAppName)
	} else {
		l.Entry.Data[entryKeyAppName] = name
	}
}

// Level return the logger level
func (l *LogrusLogger) Level() Level {
	return Level(l.Logger.Level)
}

// SetLevel sets the logging level
func (l *LogrusLogger) SetLevel(level Level) {
	l.Logger.Level = logrus.Level(level)
}

// WithFields adds fields to the log and returns the logger
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := newLogrusLogger(l.cfg).(*LogrusLogger)
	newLogger.Entry = l.Entry.WithFields(fields)

	return newLogger
}
