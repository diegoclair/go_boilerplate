package logger

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/diegoclair/go_boilerplate/util/config"
	"github.com/labstack/gommon/color"
	"github.com/sirupsen/logrus"
)

// customJSONFormatter formats output in JSON with custom formats
type customJSONFormatter struct {
	logrus.JSONFormatter
	cfg config.LogConfig
}

// Format formats the message output with function name and coloring level fields if output is not in a file
func (formatter *customJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	resultSlice, err := formatter.JSONFormatter.Format(entry)
	if err != nil {
		return nil, err
	}

	result := string(resultSlice)

	msg := getLogFuncName() + ": " + entry.Message
	result = strings.Replace(result, `"msg":"`+entry.Message+`"`, `"msg":"`+msg+`"`, 1)

	fileAndLine := getLogFilename() + ":" + strconv.Itoa(getLogFileLine())
	result = strings.Replace(result, `"file":""`, `"file":"`+fileAndLine+`"`, 1)

	if !formatter.cfg.LogToFile {
		level := entry.Level.String()
		levelUpper := strings.ToUpper(level)
		levelColor := level
		clr := color.New()
		clr.Enable()

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
	}
	resultSlice = []byte(result)

	return resultSlice, nil
}

func getLogFuncName() (funcName string) {
	funcName, _, _ = getRuntimeData()
	return
}

func getLogFilename() (filename string) {
	_, filename, _ = getRuntimeData()
	return
}

func getLogFileLine() (line int) {
	_, _, line = getRuntimeData()
	return
}

func getRuntimeData() (funcname, filename string, line int) {

	pc, filePath, line, ok := runtime.Caller(7)
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
