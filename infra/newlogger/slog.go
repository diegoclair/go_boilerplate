package newlogger

import (
	"os"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"golang.org/x/exp/slog"
)

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
	logger.Logger = slog.New(newCustomJSONFormatter(os.Stdout, opts, cfg.Log))
	return logger
}
