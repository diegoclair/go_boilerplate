package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/diegoclair/go_boilerplate/infra/config"
	"github.com/gookit/color"
	"github.com/labstack/gommon/log"

	"log/slog"
)

type customJSONFormatter struct {
	slog.Handler
	w    io.Writer
	cfg  config.LogConfig
	attr []slog.Attr
}

func newCustomJSONFormatter(w io.Writer, opts slog.HandlerOptions, cfg config.Config) *customJSONFormatter {

	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("error obtaining host name: %v", err)
	}
	return &customJSONFormatter{
		Handler: slog.NewJSONHandler(w, &opts),
		w:       w,
		cfg:     cfg.Log,
		attr: []slog.Attr{
			slog.String("app", cfg.App.Name),
			slog.String("hostname", hostname),
		},
	}
}

func (f *customJSONFormatter) Handle(ctx context.Context, r slog.Record) error {

	funcName, fileName, fileLine := f.getRuntimeData()

	level := f.getLevel(r.Level)

	buf := strings.Builder{}
	buf.WriteByte('{')
	buf.WriteString(fmt.Sprintf(`"time":"%s"`, r.Time.Format("2006-01-02T15:04:05")))
	buf.WriteString(fmt.Sprintf(`,"level":"%s"`, level))
	buf.WriteString(fmt.Sprintf(`,"file":"%s:%d"`, fileName, fileLine))
	buf.WriteString(fmt.Sprintf(`,"msg":"%s: %s"`, funcName, r.Message))

	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(fmt.Sprintf(`,"%s":"%s"`, a.Key, a.Value.Any()))
		return true
	})
	for _, attr := range f.attr {
		buf.WriteString(fmt.Sprintf(`,"%s":"%s"`, attr.Key, attr.Value.Any()))
	}
	buf.WriteByte('}')

	_, err := fmt.Fprintln(f.w, f.applyLevelColor(buf.String(), level))
	if err != nil {
		return err
	}

	return nil
}

func (f *customJSONFormatter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return f.Handler.WithAttrs(f.attr)
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
		case slog.LevelInfo.String():
			levelColor = color.Blue.Render(levelUpper)
		case slog.LevelDebug.String():
			levelColor = color.Magenta.Render(levelUpper)
		case slog.LevelWarn.String():
			levelColor = color.Yellow.Render(levelUpper)
		case slog.LevelError.String():
			levelColor = color.Red.Render(levelUpper)
		case LevelFatal:
			levelColor = color.Bold.Render(color.Red.Render(levelUpper))

		default:
			levelColor = levelUpper
		}

		return strings.Replace(fullMsg, `"level":"`+level+`"`, `"level":"`+levelColor+`"`, 1)
	}
	return fullMsg
}

func (f *customJSONFormatter) getLevel(level slog.Level) string {
	if l, ok := CustomLevels[int(level)]; ok {
		return l
	}
	return level.String()
}

func (f *customJSONFormatter) getRuntimeData() (funcname, filename string, line int) {
	pc, filePath, line, ok := runtime.Caller(5)
	if !ok {
		panic("Could not get context info for logger!")
	}
	filename = filepath.Base(filePath)
	funcPath := runtime.FuncForPC(pc).Name()
	funcname = funcPath[strings.LastIndex(funcPath, ".")+1:]

	//handle go func called inside of a function
	/*
		for example, we have a func Example() and inside of it, we have a go func() without a name, the it will output funcname as func1, with this handle, it will
		output func name as Example.func1
	*/
	if strings.Contains(funcname, "func") {
		funcBefore := funcPath[:strings.LastIndex(funcPath, ".")]
		funcname = funcPath[strings.LastIndex(funcBefore, ".")+1:]
	}
	return
}
