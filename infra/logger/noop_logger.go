package logger

import "context"

func NewNoop() Logger {
	return &NoopLogger{}
}

type NoopLogger struct{}

func (l *NoopLogger) Info(ctx context.Context, msg string)                        {}
func (l *NoopLogger) Infof(ctx context.Context, msg string, args ...any)          {}
func (l *NoopLogger) Infow(ctx context.Context, msg string, keyAndValues ...any)  {}
func (l *NoopLogger) Debug(ctx context.Context, msg string)                       {}
func (l *NoopLogger) Debugf(ctx context.Context, msg string, args ...any)         {}
func (l *NoopLogger) Debugw(ctx context.Context, msg string, keyAndValues ...any) {}
func (l *NoopLogger) Warn(ctx context.Context, msg string)                        {}
func (l *NoopLogger) Warnf(ctx context.Context, msg string, args ...any)          {}
func (l *NoopLogger) Warnw(ctx context.Context, msg string, keyAndValues ...any)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string)                       {}
func (l *NoopLogger) Errorf(ctx context.Context, msg string, args ...any)         {}
func (l *NoopLogger) Errorw(ctx context.Context, msg string, keyAndValues ...any) {}
func (l *NoopLogger) Fatal(ctx context.Context, msg string)                       {}
func (l *NoopLogger) Fatalf(ctx context.Context, msg string, args ...any)         {}
func (l *NoopLogger) Fatalw(ctx context.Context, msg string, keyAndValues ...any) {}
func (l *NoopLogger) Print(args ...any)                                           {}
