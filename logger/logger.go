package logger

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Info(format string, v ...any) {
	log.Info().Msgf(format, v...)
}

func Warn(format string, v ...any) {
	log.Warn().Msgf(format, v...)
}

func Error(format string, v ...any) {
	log.Error().Msgf(format, v...)
}

func Debug(format string, v ...any) {
	log.Debug().Msgf(format, v...)
}

func Assert(format string, v ...any) {
	log.Printf(format, v...)
}

var CustomizeFiles []string

var loggerInfo zerolog.Logger
var loggerDebug zerolog.Logger
var loggerWarn zerolog.Logger
var loggerError zerolog.Logger
var loggerAssert zerolog.Logger
var loggerTrace zerolog.Logger

func InitLog() {
	//日志级别设置，默认Info
	level := config.GetValueStringDefault("server.logger.level", "info")
	if l, err := zerolog.ParseLevel(strings.ToLower(level)); err != nil {
		log.Warn().Msgf("日志设置异常，将使用默认级别 INFO")
	} else {
		zerolog.SetGlobalLevel(l)
	}
	//时间格式设置
	timeFieldFormat := config.GetValueStringDefault("server.logger.time.format", time.RFC3339)
	zerolog.TimeFieldFormat = timeFieldFormat
	//设置日志输出
	out := zerolog.ConsoleWriter{Out: os.Stderr}
	colorEnable := config.GetValueBoolDefault("server.logger.color.enable", false)
	if colorEnable {
		out = zerolog.ConsoleWriter{Out: colorable.NewColorableStdout()}
	}
	out.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" [%s] [%-2s]", config.GetValueStringDefault("base.application.name", "isc-gobase"), i))
	}
	log.Logger = log.Logger.Output(out).With().Caller().Timestamp().Logger()

	//添加hook
	levelInfoHook := zerolog.HookFunc(func(e *zerolog.Event, l zerolog.Level, msg string) {
		//levelName := l.String()
		e1 := e
		if isc.ListAny[string](CustomizeFiles, func(t string) bool {
			return zerolog.CallerFieldName == t
		}) {
			//日志修改日志级别为debug并输出日志
			logger := log.With().Logger()
			logger.Debug().Msg(msg)
		}

		switch l {
		case zerolog.DebugLevel:
			e1 = loggerDebug.Debug()
		case zerolog.InfoLevel:
			e1 = loggerInfo.Info()
		case zerolog.WarnLevel:
			e1 = loggerWarn.Warn()
		case zerolog.ErrorLevel:
			e1 = loggerError.Error()
		case zerolog.TraceLevel:
			e1 = loggerTrace.Trace()
		default:
			e1 = loggerAssert.Log()
		}
		e1.Msg(msg)
	})
	log.Hook(levelInfoHook)
}

func initLoggerFile(logDir string, fileName string) zerolog.Logger {
	var l zerolog.Logger
	logFile := filepath.Join(logDir, fileName)
	if file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err == nil {
		l = log.With().Logger()
		l.Output(file)
	}
	return l
}

func init() {
	// 创建日志目录
	logDir := filepath.Join(".", "logs")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.Mkdir(logDir, os.ModePerm)
	}
	// 创建日志文件
	loggerInfo = initLoggerFile(logDir, "app-info.log")
	loggerDebug = initLoggerFile(logDir, "app-debug.log")
	loggerWarn = initLoggerFile(logDir, "app-warn.log")
	loggerError = initLoggerFile(logDir, "app-error.log")
	loggerAssert = initLoggerFile(logDir, "app-assert.log")

}
