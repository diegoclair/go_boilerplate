package newlogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/gookit/color"

	"golang.org/x/exp/slog"
)

type customJSONFormatter struct {
	slog.Handler
	w   io.Writer
	cfg config.LogConfig
}

func newCustomJSONFormatter(w io.Writer, opts slog.HandlerOptions, cfg config.LogConfig) *customJSONFormatter {
	return &customJSONFormatter{
		Handler: slog.NewJSONHandler(w, &opts),
		w:       w,
		cfg:     cfg,
	}
}

func (f *customJSONFormatter) Handle(ctx context.Context, r slog.Record) error {

	level := r.Level.String()
	fields := make(map[string]any, r.NumAttrs())

	fields["level"] = level
	fields["time"] = r.Time.Format("2006-01-02T15:04:05")

	frames := runtime.CallersFrames([]uintptr{r.PC})
	frame, _ := frames.Next()

	fields["file"] = filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line)

	//[0] is the package of caller, [1] is the caller func name
	funcName := strings.Split(frame.Function, ".")[1]
	fields["msg"] = funcName + ": " + r.Message

	b, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(f.w, f.applyLevelColor(string(b), level))
	if err != nil {
		return err
	}

	return nil
}

func (f *customJSONFormatter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return f.Handler.WithAttrs(attrs)
}

func (f *customJSONFormatter) WithGroup(name string) slog.Handler {
	return f.Handler.WithGroup(name)
}

func (f *customJSONFormatter) Enabled(ctx context.Context, level slog.Level) bool {
	return f.Handler.Enabled(ctx, level)
}

func (f *customJSONFormatter) applyLevelColor(fullMsg, level string) string {

	if !f.cfg.LogToFile {
		level := level
		levelUpper := strings.ToUpper(level)
		levelColor := ""

		switch level {
		case slog.LevelDebug.String():
			levelColor = color.Magenta.Render(levelUpper)
		case slog.LevelInfo.String():
			levelColor = color.Blue.Render(levelUpper)
		case slog.LevelWarn.String():
			levelColor = color.Yellow.Render(levelUpper)
		case slog.LevelError.String():
			levelColor = color.Red.Render(levelUpper)

		default:
			levelColor = levelUpper
		}

		return strings.Replace(fullMsg, `"level":"`+level+`"`, `"level":"`+levelColor+`"`, 1)
	}
	return fullMsg
}
