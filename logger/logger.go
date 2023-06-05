package logger

import (
	"bytes"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/store"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
	"sync"
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
//var loggerMap map[string]*logrus.Logger
var loggerMap cmap.ConcurrentMap
//var rotateMap map[string]*rotatelogs.RotateLogs
var rotateMap cmap.ConcurrentMap
var rootLogger *logrus.Logger

func init() {
	loggerMap = cmap.New()
	rotateMap = cmap.New()
	rootLogger = Group("root")

	_gColor := config.GetValueBoolDefault("base.logger.color.enable", false)
	gColor = _gColor
}

func Group(groupNames... string) *logrus.Logger {
	var resultLogger *logrus.Logger
	groupNamesOfUnContain := []string{}
	for _, groupName := range groupNames {
		if logger, exit := loggerMap.Get(groupName); exit {
			resultLogger = logger.(*logrus.Logger)
		} else {
			groupNamesOfUnContain = append(groupNamesOfUnContain, groupName)
		}
	}

	if resultLogger != nil {
		return resultLogger
	}  else {
		resultLogger = logrus.New()
		resultLogger.SetReportCaller(true)
		formatters := &StandardFormatter{}
		resultLogger.Formatter = formatters

		loggerDir := config.GetValueStringDefault("base.logger.home", "./logs/")
		resultLogger.AddHook(lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: rotateLogWithCache(loggerDir, "debug"),
			logrus.InfoLevel:  rotateLogWithCache(loggerDir, "info"),
			logrus.WarnLevel:  rotateLogWithCache(loggerDir, "warn"),
			logrus.ErrorLevel: rotateLogWithCache(loggerDir, "error"),
			logrus.PanicLevel: rotateLogWithCache(loggerDir, "panic"),
			logrus.FatalLevel: rotateLogWithCache(loggerDir, "fatal"),
		}, formatters))
	}

	// 值最大的级别，对应的level最小，比如Debug对应的数值比Info要大
	maxValueLevel := logrus.PanicLevel
	for _, groupName := range groupNamesOfUnContain {
		var finalGroupLevel string
		rootLevel := config.GetValueStringDefault("base.logger.level", "info")
		groupLevel := config.GetValueString("base.logger.group." + groupName + ".level")
		if groupLevel != "" {
			finalGroupLevel = groupLevel
		} else {
			finalGroupLevel = rootLevel
		}

		lgLevel, err := logrus.ParseLevel(finalGroupLevel)
		if err != nil {
			lgLevel = logrus.InfoLevel
		}

		if lgLevel > maxValueLevel {
			maxValueLevel = lgLevel
		}
	}

	resultLogger.SetLevel(maxValueLevel)

	for _, groupName := range groupNamesOfUnContain {
		loggerMap.Set(groupName, resultLogger)
	}
	return resultLogger
}

func doGroup(groupName string) *logrus.Logger {
	if groupName == "" {
		return rootLogger
	}
	if logger, exit := loggerMap.Get(groupName); exit {
		return logger.(*logrus.Logger)
	}

	if loggerMap == nil {
		loggerMap = cmap.New()
	}
	logger := logrus.New()
	logger.SetReportCaller(true)
	formatters := &StandardFormatter{}
	logger.Formatter = formatters

	loggerDir := config.GetValueStringDefault("base.logger.home", "./logs/")
	logger.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotateLogWithCache(loggerDir, "debug"),
		logrus.InfoLevel:  rotateLogWithCache(loggerDir, "info"),
		logrus.WarnLevel:  rotateLogWithCache(loggerDir, "warn"),
		logrus.ErrorLevel: rotateLogWithCache(loggerDir, "error"),
		logrus.PanicLevel: rotateLogWithCache(loggerDir, "panic"),
		logrus.FatalLevel: rotateLogWithCache(loggerDir, "fatal"),
	}, formatters))

	var finalGroupLevel string
	rootLevel := config.GetValueStringDefault("base.logger.level", "info")
	groupLevel := config.GetValueString("base.logger.group." + groupName + ".level")
	if groupLevel != "" {
		finalGroupLevel = groupLevel
	} else {
		finalGroupLevel = rootLevel
	}

	lgLevel, err := logrus.ParseLevel(finalGroupLevel)
	if err != nil {
		lgLevel = logrus.InfoLevel
	}
	logger.SetLevel(lgLevel)

	loggerMap.Set(groupName, logger)
	return logger
}

func InitLog() {
	rootLogger = Group("root")
	loggerDir := config.GetValueStringDefault("base.logger.home", "./logs/")
	rootLogger.AddHook(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotateLog(loggerDir, "debug"),
		logrus.InfoLevel:  rotateLog(loggerDir, "info"),
		logrus.WarnLevel:  rotateLog(loggerDir, "warn"),
		logrus.ErrorLevel: rotateLog(loggerDir, "error"),
		logrus.PanicLevel: rotateLog(loggerDir, "panic"),
		logrus.FatalLevel: rotateLog(loggerDir, "fatal"),
	}, &StandardFormatter{}))
	lgLevel, err := logrus.ParseLevel(config.GetValueStringDefault("base.logger.level", "info"))
	if err != nil {
		lgLevel = logrus.InfoLevel
	}
	rootLogger.SetLevel(lgLevel)

	_gColor := config.GetValueBoolDefault("base.logger.color.enable", false)
	gColor = _gColor

	listener.AddListener(listener.EventOfConfigChange, ConfigChangeListener)
}

func ConfigChangeListener(event listener.BaseEvent) {
	ev := event.(listener.ConfigChangeEvent)
	if ev.Key == "base.logger.level" {
		SetGlobalLevel(ev.Value)
	} else if strings.HasPrefix(ev.Key, "base.logger.group") {
		words := strings.Split(ev.Key, ".")
		if len(words) != 5 {
			return
		}
		_group := words[3]
		_level := ev.Value
		le, err := logrus.ParseLevel(_level)
		if err != nil {
			return
		}
		Group(_group).SetLevel(le)
	}
}

func SetGlobalLevel(strLevel string) {
	level, err := logrus.ParseLevel(strLevel)
	if err == nil {
		rootLogger.SetLevel(level)
	}
}

func InfoDirect(v ...any) {
	rootLogger.Info(v...)
}

func WarnDirect(v ...any) {
	rootLogger.Warn(v...)
}

func ErrorDirect(v ...any) {
	rootLogger.Error(v...)
}

func FatalDirect(v ...any) {
	rootLogger.Fatal(v...)
}

func PanicDirect(v ...any) {
	rootLogger.Panic(v...)
}

func DebugDirect(v ...any) {
	rootLogger.Debug(v...)
}

func TraceDirect(v ...any) {
	rootLogger.Trace(v...)
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

func Trace(format string, v ...any) {
	rootLogger.Tracef(format, v...)
}

func Panic(format string, v ...any) {
	rootLogger.Panicf(format, v...)
}

func Fatal(format string, v ...any) {
	rootLogger.Fatalf(format, v...)
}

func Record(level, format string, v ...any) {
	switch strings.ToLower(level) {
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
	if rotateMap == nil {
		rotateMap = cmap.New()
	}

	if path == "" {
		path = "./logs/"
	}

	maxSizeStr := config.GetValueStringDefault("base.logger.rotate.max-size", "300MB")
	maxHistoryStr := config.GetValueStringDefault("base.logger.rotate.max-history", "60d")
	rotateTimeStr := config.GetValueStringDefault("base.logger.rotate.time", "1d")

	rotateOptions := []rotatelogs.Option{rotatelogs.WithLinkName(path + "app-" + level + ".log")}
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

	data, _ := rotatelogs.New(path+"app-"+level+".%Y%m%d.log", rotateOptions...)
	rotateMap.Set(path+"-"+level, data)
	return data
}

func rotateLogWithCache(path, level string) *rotatelogs.RotateLogs {
	if pRotateValue, exist := rotateMap.Get(path+"-"+level); exist {
		return pRotateValue.(*rotatelogs.RotateLogs)
	}

	return rotateLog(path, level)
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
		frame := getCallerFrame()
		funPath = fmt.Sprintf("%s:%d#%s", shortLogPath(frame.File), frame.Line, functionName(frame))
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
			levelColor = green
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel:
			levelColor = red
		case logrus.FatalLevel:
			levelColor = red
		case logrus.PanicLevel:
			levelColor = red
		}
		newLog = fmt.Sprintf("[%s] \x1b[%dm%s [%s]\x1b[0m [%s] [%v] \x1b[%dm%s\x1b[0m \x1b[%dm%s\x1b[0m %s %s\n",
			timestamp,
			black,
			os.Getenv("HOSTNAME"),
			config.GetValueStringDefault("base.application.name", "isc-gobase"),
			store.Get(constants.TRACE_HEAD_ID), store.Get(constants.TRACE_HEAD_USER_ID),
			levelColor,
			strings.ToUpper(entry.Level.String()),
			black,
			funPath,
			entry.Message,
			fieldsStr)
	} else {
		newLog = fmt.Sprintf("[%s] %s [%s] [%s] [%v] %s %s %s %s\n",
			timestamp,
			os.Getenv("HOSTNAME"),
			config.GetValueStringDefault("base.application.name", "isc-gobase"),
			store.Get(constants.TRACE_HEAD_ID), store.Get(constants.TRACE_HEAD_USER_ID),
			strings.ToUpper(entry.Level.String()),
			funPath,
			entry.Message,
			fieldsStr)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

const (
	maximumCallerDepth    int = 25
	knownBaseLoggerFrames int = 5
)

var callerInitOnce sync.Once
var minimumCallerDepth = 0
var baseLoggerPackage string

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func getCallerFrame() *runtime.Frame {
	pcs := make([]uintptr, maximumCallerDepth)
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "logger.getCallerFrame") {
				baseLoggerPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownBaseLoggerFrames
	})

	pcs = make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if pkg != baseLoggerPackage && pkg != "github.com/sirupsen/logrus" {
			return &f
		}
	}
	return nil
}

func functionName(frame *runtime.Frame) string {
	pathMeta := strings.Split(frame.Function, ".")
	if len(pathMeta) > 1 {
		return pathMeta[len(pathMeta)-1]
	}
	return frame.Function
}

func shortLogPath(logPath string) string {
	loggerPath := config.GetValueStringDefault("base.logger.path.type", "short")
	if loggerPath == "short" {
		pathMeta := strings.Split(logPath, string(os.PathSeparator))
		if len(pathMeta) > 3 {
			return pathMeta[len(pathMeta)-3] + string(os.PathSeparator) + pathMeta[len(pathMeta)-2] + string(os.PathSeparator) + pathMeta[len(pathMeta)-1]
		}
		return logPath
	} else if loggerPath == "full" {
		pathMeta := strings.Split(logPath, "@2/project")
		if len(pathMeta) > 1 {
			pathMeta[0] = "../.."
			return strings.Join(pathMeta, "")
		}
		return logPath
	} else {
		return logPath
	}
}
