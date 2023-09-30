package logger

import (
	"context"

	"github.com/diegoclair/go_boilerplate/infra/config"
)

const (
	AccountUUIDKey = "account_uuid"
	ErrorKey       = "error"
)

type Logger interface {
	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, msg string, args ...any)
	Infow(ctx context.Context, msg string, keyAndValues ...any)
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, msg string, args ...any)
	Debugw(ctx context.Context, msg string, keyAndValues ...any)
	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, msg string, args ...any)
	Warnw(ctx context.Context, msg string, keyAndValues ...any)
	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, msg string, args ...any)
	Errorw(ctx context.Context, msg string, keyAndValues ...any)
	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, msg string, args ...any)
	Fatalfw(ctx context.Context, msg string, keyAndValues ...any)
	Print(args ...any)
}

func New(cfg config.Config) Logger {
	return newSlogLogger(cfg)
}
