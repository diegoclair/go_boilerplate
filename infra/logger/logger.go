package logger

import (
	"context"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/diegoclair/go_utils/logger"
)

func New(cfg *config.Config) logger.Logger {

	params := logger.LogParams{
		AppName:                  cfg.App.Name,
		DebugLevel:               cfg.Log.Debug,
		AddAttributesFromContext: addDefaultAttributesToLogger,
	}

	return logger.New(params)
}

func addDefaultAttributesToLogger(ctx context.Context) []any {
	args := []any{}
	if sessionCode, ok := getSession(ctx); ok {
		args = append(args, "session", sessionCode)
	}

	if accountUUID, ok := getAccountUUID(ctx); ok {
		args = append(args, "account_uuid", accountUUID)
	}

	return args
}

func getContextValue(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}

	value := ctx.Value(key)
	if value == nil {
		return ""
	}

	return value.(string)
}

func getSession(ctx context.Context) (string, bool) {
	sessionCode := getContextValue(ctx, string(infra.SessionKey))
	if sessionCode == "" {
		return "", false
	}

	return sessionCode, true
}

func getAccountUUID(ctx context.Context) (string, bool) {
	accountUUID := getContextValue(ctx, string(infra.AccountUUIDKey))
	if accountUUID == "" {
		return "", false
	}

	return accountUUID, true
}
