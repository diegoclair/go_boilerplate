package logger

import (
	"context"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/logger"
)

func NewLogger(appName string, debugLevel bool) logger.Logger {
	params := logger.Params{
		AppName:          appName,
		DebugLevel:       debugLevel,
		ContextExtractor: addDefaultAttributesToLogger,
	}
	return logger.New(params)
}

func addDefaultAttributesToLogger(ctx context.Context) []logger.Field {
	args := []logger.Field{}

	if sessionCode, ok := getContextValue[string](ctx, infra.SessionKey); ok {
		args = append(args, logger.Attr("session", sessionCode))
	}

	if accountUUID, ok := getContextValue[string](ctx, infra.AccountUUIDKey); ok {
		args = append(args, logger.Attr("account_uuid", accountUUID))
	}

	return args
}

func getContextValue[T comparable](ctx context.Context, key infra.Key) (T, bool) {
	var zero T
	if ctx == nil {
		return zero, false
	}

	value := ctx.Value(key)
	if value == nil {
		return zero, false
	}

	v, ok := value.(T)
	return v, ok
}
