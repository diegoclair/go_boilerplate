package logger

import "github.com/diegoclair/go_boilerplate/infra/config"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatalf(msg string, args ...any)
	Fatal(msg string, args ...any)
	Print(args ...any)
}

func New(cfg config.Config) Logger {
	return newSlogLogger(cfg)
}
