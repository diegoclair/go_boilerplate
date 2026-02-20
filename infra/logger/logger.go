package logger

import (
	"context"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_utils/logger"
)

func NewLogger(appName string, debugLevel bool) logger.Logger {
	params := logger.LogParams{
		AppName:                  appName,
		DebugLevel:               debugLevel,
		AddAttributesFromContext: addDefaultAttributesToLogger,
	}
	return logger.New(params)
}

func addDefaultAttributesToLogger(ctx context.Context) []logger.LogField {
	args := []logger.LogField{}

	if sessionCode, ok := getContextValue[string](ctx, infra.SessionKey); ok {
		args = append(args, logger.String("session", sessionCode))
	}

	if accountUUID, ok := getContextValue[string](ctx, infra.AccountUUIDKey); ok {
		args = append(args, logger.String("account_uuid", accountUUID))
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
