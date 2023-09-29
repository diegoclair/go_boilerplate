package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"golang.org/x/exp/slog"
)

const (
	LevelFatal     = "FATAL"
	LevelFatalCode = 60
)

var CustomLevels = map[int]string{
	LevelFatalCode: LevelFatal, //high number to avoid conflict with slog levels
}

type SlogLogger struct {
	cfg config.Config
	*slog.Logger
}

func newSlogLogger(cfg config.Config) *SlogLogger {

	logger := &SlogLogger{cfg: cfg}

	opts := slog.HandlerOptions{}

	if cfg.Log.Debug {
		opts.Level = slog.LevelDebug
	}
	logger.Logger = slog.New(newCustomJSONFormatter(os.Stdout, opts, cfg))
	return logger
}

func (l *SlogLogger) Fatalf(msg string, args ...any) {
	l.Logger.Log(context.TODO(), LevelFatalCode, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (l *SlogLogger) Fatal(msg string, args ...any) {
	l.Logger.Log(context.TODO(), LevelFatalCode, msg, args...)
	os.Exit(1)
}

func (l *SlogLogger) Print(args ...any) {
	l.Logger.Log(context.TODO(), 0, "", args...)
}
