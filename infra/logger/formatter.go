package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/diegoclair/go-boilerplate/util/config"
	"github.com/labstack/gommon/color"
	"github.com/sirupsen/logrus"
)

// customJSONFormatter formats output in JSON with custom formats
type customJSONFormatter struct {
	logrus.JSONFormatter
	cfg config.LogConfig
}

// Format formats the log output, coloring status and level fields if output is not in a file
func (formatter *customJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	resultSlice, err := formatter.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	result := string(resultSlice)

	msg := formatter.getFuncName() + ": " + entry.Message
	result = strings.Replace(result, `"msg":"`+entry.Message+`"`, `"msg":"`+msg+`"`, 1)

	if !formatter.cfg.LogToFile {

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
	}
	resultSlice = []byte(result)

	return resultSlice, nil
}

func (formatter *customJSONFormatter) getFuncName() string {
	pc, _, _, ok := runtime.Caller(6)
	if !ok {
		panic("Could not get context info for logger!")
	}
	funcPath := runtime.FuncForPC(pc).Name()

	funcName := funcPath[strings.LastIndex(funcPath, ".")+1:]

	//handle go func called inside of a function
	/*
		for example, we have a func Example() and inside of it, we have a go func() without a name, the it will output funcname as func1, with this handle, it will
		output func name as Example.func1
	*/
	if strings.Contains(funcName, "func") {
		funcBefore := funcPath[:strings.LastIndex(funcPath, ".")]
		funcName = funcPath[strings.LastIndex(funcBefore, ".")+1:]
	}
	return funcName
}
