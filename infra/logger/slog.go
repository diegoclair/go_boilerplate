package logger

import (
	"context"
	"fmt"
	"os"

	"log/slog"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/config"
)

const sessionCodeKey = "session_code"

const (
	LevelFatal     = "FATAL"
	LevelFatalCode = 60
)

var CustomLevels = map[int]string{
	LevelFatalCode: LevelFatal, //high number to avoid conflict with slog levels
}

type SlogLogger struct {
	cfg *config.Config
	*slog.Logger
}

func newSlogLogger(cfg *config.Config) *SlogLogger {
	logger := &SlogLogger{cfg: cfg}
	opts := slog.HandlerOptions{}
	if cfg.Log.Debug {
		opts.Level = slog.LevelDebug
	}
	logger.Logger = slog.New(newCustomJSONFormatter(os.Stdout, opts, cfg))
	return logger
}

func (l *SlogLogger) Info(ctx context.Context, msg string) {
	l.Logger.InfoContext(ctx, msg, l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Infof(ctx context.Context, msg string, args ...any) {
	l.Logger.InfoContext(ctx, fmt.Sprintf(msg, args...), l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Infow(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.InfoContext(ctx, msg, append(l.withDefaultArgs(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Debug(ctx context.Context, msg string) {
	l.Logger.DebugContext(ctx, msg, l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Debugf(ctx context.Context, msg string, args ...any) {
	l.Logger.DebugContext(ctx, fmt.Sprintf(msg, args...), l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Debugw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.DebugContext(ctx, msg, append(l.withDefaultArgs(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string) {
	l.Logger.WarnContext(ctx, msg, l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Warnf(ctx context.Context, msg string, args ...any) {
	l.Logger.WarnContext(ctx, fmt.Sprintf(msg, args...), l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Warnw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.WarnContext(ctx, msg, append(l.withDefaultArgs(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string) {
	l.Logger.ErrorContext(ctx, msg, l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Errorf(ctx context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(ctx, fmt.Sprintf(msg, args...), l.withDefaultArgs(ctx)...)
}

func (l *SlogLogger) Errorw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.ErrorContext(ctx, msg, append(l.withDefaultArgs(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Fatal(ctx context.Context, msg string) {
	l.Logger.Log(ctx, LevelFatalCode, msg, l.withDefaultArgs(ctx)...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.Logger.Log(ctx, LevelFatalCode, fmt.Sprintf(msg, args...), l.withDefaultArgs(ctx)...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.Log(ctx, LevelFatalCode, msg, append(l.withDefaultArgs(ctx), keyAndValues...)...)
	os.Exit(1)
}

func (l *SlogLogger) Print(args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, "", args...)
}

func (l *SlogLogger) getContextValue(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}

	value := ctx.Value(key)
	if value == nil {
		return ""
	}

	return value.(string)
}

func (l *SlogLogger) getSession(ctx context.Context) (string, bool) {
	sessionCode := l.getContextValue(ctx, string(infra.SessionKey))
	if sessionCode == "" {
		return "", false
	}

	return sessionCode, true
}

func (l *SlogLogger) getAccountUUID(ctx context.Context) (string, bool) {
	accountUUID := l.getContextValue(ctx, string(infra.AccountUUIDKey))
	if accountUUID == "" {
		return "", false
	}

	return accountUUID, true
}

func (l *SlogLogger) withDefaultArgs(ctx context.Context) []any {
	args := []any{}
	if sessionCode, ok := l.getSession(ctx); ok {
		args = append(args, sessionCodeKey, sessionCode)
	}

	if accountUUID, ok := l.getAccountUUID(ctx); ok {
		args = append(args, "account_uuid", accountUUID)
	}

	return args
}
