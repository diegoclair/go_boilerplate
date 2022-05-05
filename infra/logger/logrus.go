package logger

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/diegoclair/go-boilerplate/infra/auth"
	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/sirupsen/logrus"
)

// entryKey standard values
const (
	entryKeyAppName  = "app"
	entryKeyRecordID = "id"
)

func newLogrusLogger(cfg config.LogConfig) Logger {
	if cfg.LogToFile {
		file, err := os.Create(cfg.Path)
		if err != nil {
			fmt.Printf("Error to create log file for library: %s\n", err.Error())
			panic(err)
		}
		logrus.SetOutput(file)
	}

	logrus.SetFormatter(&coloredJSONFormatter{cfg: cfg})

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Errorf("Error obtaining host name: %v", err)
	}

	entry := logrus.WithFields(logrus.Fields{
		"hostname": hostname,
	})

	return &LogrusLogger{cfg, entry}
}

// LogrusLogger is the default application logger
type LogrusLogger struct {
	cfg config.LogConfig
	*logrus.Entry
}

// NewSessionLogger returns an instance of log with session code field
func (l *LogrusLogger) NewSessionLogger(ctx context.Context) (context.Context, Logger) {

	var instance Logger
	sessionCode := ctx.Value(auth.SessionKey)
	if sessionCode == nil {
		instance = newLogrusLogger(l.cfg)
		return ctx, instance
	}

	var sessionCodeKey = "logger-" + sessionCode.(string)
	vl := ctx.Value(sessionCodeKey)
	if vl != nil {
		return ctx, vl.(Logger)
	}
	instance = newLogrusLogger(l.cfg)
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

// Output returns the echo logging output
func (l *LogrusLogger) Output() io.Writer {
	return l.Logger.Writer()
}

// SetOutput sets the echo logging output
func (l *LogrusLogger) SetOutput(output io.Writer) {
	l.Logger.Out = output
}

// RecordID return the entry ID
func (l *LogrusLogger) RecordID() string {
	recordID, ok := l.Entry.Data[entryKeyRecordID]
	if !ok {
		return ""
	}

	return recordID.(string)
}

// SetRecordID sets the entry ID
func (l *LogrusLogger) SetRecordID(id string) {
	if id == "" {
		delete(l.Entry.Data, entryKeyRecordID)
	} else {
		l.Entry.Data[entryKeyRecordID] = id
	}
}

// InfoWriter returns the io.Writer for info level
func (l *LogrusLogger) InfoWriter() io.Writer {
	return l.WriterLevel(logrus.InfoLevel)
}

// ErrorWriter returns the io.Writer for error level
func (l *LogrusLogger) ErrorWriter() io.Writer {
	return l.WriterLevel(logrus.ErrorLevel)
}

// FatalWriter returns the io.Writer for fatal level
func (l *LogrusLogger) FatalWriter() io.Writer {
	return l.WriterLevel(logrus.FatalLevel)
}

// WithFields adds fields to the log and returns the logger
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	newLogger := newLogrusLogger(l.cfg).(*LogrusLogger)
	newLogger.Entry = l.Entry.WithFields(fields)

	return newLogger
}

// Fields returns the logger fields
func (l *LogrusLogger) Fields() map[string]interface{} {
	return l.Entry.Data
}
