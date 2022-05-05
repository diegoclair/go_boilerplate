package logger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/labstack/gommon/color"
	"github.com/sirupsen/logrus"
)

// coloredJSONFormatter formats output in JSON with colored strings
type coloredJSONFormatter struct {
	logrus.JSONFormatter
	cfg config.LogConfig
}

// Format formats the log output, coloring status and level fields if output is not in a file
func (formatter *coloredJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	resultSlice, err := formatter.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	if !formatter.cfg.LogToFile {
		result := string(resultSlice)

		level := entry.Level.String()
		levelUpper := strings.ToUpper(level)
		levelColor := level
		clr := color.New()

		switch level {
		case logrus.DebugLevel.String():
			levelColor = clr.Magenta(levelUpper)
		case logrus.InfoLevel.String():
			levelColor = clr.Blue(levelUpper)
		case logrus.WarnLevel.String():
			levelColor = clr.Yellow(levelUpper)
		case logrus.ErrorLevel.String():
			levelColor = clr.Red(levelUpper)
		case logrus.FatalLevel.String():
			levelColor = clr.Bold(clr.Red(levelUpper))
		default:
			levelColor = levelUpper
		}

		result = strings.Replace(result, `"level":"`+level+`"`, `"level":"`+levelColor+`"`, 1)

		status, ok := entry.Data["status"].(int)
		if ok {
			statusColor := fmt.Sprint(status)

			if status >= 500 {
				statusColor = clr.Red(status)
			} else if status >= 400 {
				statusColor = clr.Yellow(status)
			} else if status >= 300 {
				statusColor = clr.Cyan(status)
			} else {
				statusColor = clr.Green(status)
			}

			result = strings.Replace(result, `"status":`+strconv.Itoa(status), `"status":`+statusColor, 1)
		}

		resultSlice = []byte(result)
	}

	return resultSlice, nil
}
