package baselog

import (
	"bytes"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"time"
)

const (
	white  = 29
	black  = 30
	red    = 31
	green  = 32
	yellow = 33
	purple = 35
	blue   = 36
	gray   = 37
)

var gColor = false
var loggerMap map[string]*logrus.Logger
var rotateMap map[string]*rotatelogs.RotateLogs
var rootLogger *logrus.Logger

func init() {
	_loggerMap := map[string]*logrus.Logger{}
	loggerMap = _loggerMap
	_rotateMap := map[string]*rotatelogs.RotateLogs{}
	rotateMap = _rotateMap
	rootLogger = GetLogger("root")

	_gColor := config.GetValueBoolDefault("base.logger.color.enable", false)
	gColor = _gColor
}

func GetLogger(loggerName string) *logrus.Logger {
	if logger, exit := loggerMap[loggerName]; exit {
		return logger
	}

	if loggerMap == nil {
		loggerMap = map[string]*logrus.Logger{}
	}
	logger := logrus.New()
	logger.SetReportCaller(true)
	formatters := &StandardFormatter{}
	logger.Formatter = formatters

	loggerDir := config.GetValueStringDefault("base.logger.dir", "./logs/")
	logger.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotateLog(loggerDir, "debug"),
		logrus.InfoLevel:  rotateLog(loggerDir, "info"),
		logrus.WarnLevel:  rotateLog(loggerDir, "warn"),
		logrus.ErrorLevel: rotateLog(loggerDir, "error"),
		logrus.PanicLevel: rotateLog(loggerDir, "panic"),
		logrus.FatalLevel: rotateLog(loggerDir, "fatal"),
	}, formatters))

	loggerMap[loggerName] = logger
	return logger
}

func Info(format string, v ...any) {
	rootLogger.Infof(format, v...)
}

func Warn(format string, v ...any) {
	rootLogger.Warnf(format, v...)
}

func Error(format string, v ...any) {
	rootLogger.Errorf(format, v...)
}

func Debug(format string, v ...any) {
	rootLogger.Debugf(format, v...)
}

func Panic(format string, v ...any) {
	rootLogger.Panicf(format, v...)
}

func Fatal(format string, v ...any) {
	rootLogger.Fatalf(format, v...)
}

func Record(level, format string, v ...any) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		Debug(format, v)
	case "info":
		Info(format, v)
	case "warn":
		Warn(format, v)
	case "error":
		Error(format, v)
	case "panic":
		Panic(format, v)
	case "fatal":
		Fatal(format, v)
	default:
		Debug(format, v)
	}
}

func rotateLog(path, level string) *rotatelogs.RotateLogs {
	if pRotateValue, exist := rotateMap[path+"-"+level]; exist {
		return pRotateValue
	}

	if rotateMap == nil {
		rotateMap = map[string]*rotatelogs.RotateLogs{}
	}

	if path == "" {
		path = "./logs/"
	}

	maxSizeStr := config.GetValueStringDefault("base.logger.rotate.max-size", "300MB")
	maxHistoryStr := config.GetValueStringDefault("base.logger.rotate.max-history", "60d")
	rotateTimeStr := config.GetValueStringDefault("base.logger.rotate.time", "1d")

	rotateOptions := []rotatelogs.Option{rotatelogs.WithLinkName(path+"app-"+level+".log")}
	if maxSizeStr != "" {
		rotateOptions = append(rotateOptions, rotatelogs.WithRotationSize(isc.ParseByteSize(maxSizeStr)))
	}

	_maxHistory, err := time.ParseDuration(maxHistoryStr)
	if err == nil {
		rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(_maxHistory))
	}

	_rotateTime, err := time.ParseDuration(rotateTimeStr)
	if err == nil {
		rotateOptions = append(rotateOptions, rotatelogs.WithRotationTime(_rotateTime))
	}

	data, _ := rotatelogs.New(path+"app-"+level+".log.%Y%m%d", rotateOptions...)
	rotateMap[path+"-"+level] = data
	return data
}

type StandardFormatter struct{}

func (m *StandardFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var fields []string
	for k, v := range entry.Data {
		fields = append(fields, fmt.Sprintf("%v=%v", k, v))
	}

	level := entry.Level
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var funPath string
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		funPath = fmt.Sprintf("%s %s:%d", entry.Caller.Function, fName, entry.Caller.Line)
	} else {
		funPath = fmt.Sprintf("%s", entry.Message)
	}

	var fieldsStr string
	if len(fields) != 0 {
		fieldsStr = fmt.Sprintf("[\x1b[%dm%s\x1b[0m]", blue, strings.Join(fields, " "))
	}
	var newLog string
	var levelColor = gray
	if gColor {
		switch level {
		case logrus.DebugLevel:
			levelColor = blue
		case logrus.InfoLevel:
			levelColor = blue
		case logrus.WarnLevel:
			levelColor = blue
		case logrus.ErrorLevel:
			levelColor = red
		case logrus.FatalLevel:
			levelColor = red
		case logrus.PanicLevel:
			levelColor = red
		}
		newLog = fmt.Sprintf("\x1b[%dm%s\t\x1b[0m%s \x1b[%dm%s\x1b[0m %s %s\n", levelColor, strings.ToUpper(entry.Level.String()), timestamp, black, funPath, entry.Message, fieldsStr)
	} else {
		newLog = fmt.Sprintf("%s\t %s %s %s %s\n", strings.ToUpper(entry.Level.String()), timestamp, funPath, entry.Message, fieldsStr)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}
