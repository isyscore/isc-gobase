package logger

//Packages logger provides a lightweight logging library dedicated to format logging.
//
// A global Logger Level default is INFO
// if you want change it,can be use for sample logging
// 		import "github.com/isyscore/isc-gobase/logger"
// 		logger.SetGlobalLevel("debug")
// 		log.INFO().Msg("my test")
//  		//Output: 2022-02-28 20:00:57.000  [ISC-GOBASE] [INFO] ../../pkg/domain/application.go:147 > my test
// if you want Log with no level and Message,can be use for sample logging
//		log.Log().Msg("")
// if you want Change Log Leven in runtime,can be use for sample logging
//		l = log.Logger.With().Logger()
//		l.INFO().Msg("")
// or use this
//		log.Logger.WithLevel(zerolog.NoLevel).Msgf("%s","mytest")
// isc-gobase logger provides Custom log
//	There is CustomizeFiles use for Custom log ,it will set level to debug
// 		if isc.ListAny[string](CustomizeFiles, func(t string) bool {
//			return zerolog.CallerFieldName == t
//		}) {
//			//日志修改日志级别为debug并输出日志
//			ll = zerolog.DebugLevel
//		}
//

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
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
	log.WithLevel(zerolog.NoLevel).Msgf(format, v)
}

var CustomizeFiles []string

var loggerInfo *zerolog.Logger
var loggerDebug *zerolog.Logger
var loggerWarn *zerolog.Logger
var loggerError *zerolog.Logger
var loggerAssert *zerolog.Logger
var loggerTrace *zerolog.Logger

// SetGlobalLevel sets the global override for log level. If this
// values is raised, all Loggers will use at least this value.
//
// To globally disable logs, set zerolog.GlobalLevel to Disabled.
func SetGlobalLevel(strLevel string) {
	level := zerolog.InfoLevel
	if strLevel != "" {
		if l, err := zerolog.ParseLevel(strings.ToLower(strLevel)); err != nil {
			log.Warn().Msgf("日志设置异常，将使用默认级别 INFO")
		} else {
			level = l
		}
	}
	zerolog.SetGlobalLevel(level)
}

//InitLog create a root logger. it will write to console and multiple file by level.
// note: default set root logger level is info
// it provides custom log with CustomizeFiles,if it match any caller's name ,log's level will be setting debug and output
func InitLog() {
	//日志级别设置，默认Info
	level := config.GetValueStringDefault("server.logger.level", "info")
	SetGlobalLevel(level)

	zerolog.CallerSkipFrameCount = 2
	//时间格式设置
	timeFieldFormat := config.GetValueStringDefault("server.logger.time.format", time.RFC3339)
	zerolog.TimeFieldFormat = timeFieldFormat
	//设置日志输出
	out := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: config.GetValueBoolDefault("server.logger.color.enable", false)}
	out.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" [%s] [%-2s]", config.GetValueStringDefault("base.application.name", "isc-gobase"), i))
	}
	log.Logger = log.Logger.Output(out).With().Caller().Timestamp().Logger()

	//添加hook
	levelInfoHook := zerolog.HookFunc(func(e *zerolog.Event, l zerolog.Level, msg string) {
		//levelName := l.String()
		ll := l
		e1 := e
		if isc.ListAny[string](CustomizeFiles, func(t string) bool {
			return zerolog.CallerFieldName == t
		}) {
			//日志修改日志级别为debug并输出日志
			ll = zerolog.DebugLevel
		}

		switch ll {
		case zerolog.DebugLevel:
			e1 = loggerDebug.Debug().Stack()
		case zerolog.InfoLevel:
			e1 = loggerInfo.Info().Stack()
		case zerolog.WarnLevel:
			e1 = loggerWarn.Warn().Stack()
		case zerolog.ErrorLevel:
			e1 = loggerError.Error().Stack()
		case zerolog.TraceLevel:
			e1 = loggerTrace.Trace().Stack()
		default:
			e1 = loggerAssert.WithLevel(l).Stack()
		}
		e1.Msg(msg)
	})
	log.Logger = log.Logger.Hook(levelInfoHook)
	initLogDir()
}

//initLoggerFile open or create and open log file
func initLoggerFile(logDir string, fileName string) *zerolog.Logger {
	var l zerolog.Logger
	logFile := filepath.Join(logDir, fileName)
	if file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err == nil {
		l = log.Logger.With().Logger()
		out := zerolog.ConsoleWriter{Out: file, TimeFormat: "2006-01-02 15:04:05.000", NoColor: true}
		out.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf(" [%s] [%-2s]", config.GetValueStringDefault("server.application.name", "isc-gobase"), i))
		}
		l = l.Output(out).With().Caller().Logger()
	}
	return &l
}

//initLogDir create log dir and file
func initLogDir() {
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
	loggerTrace = initLoggerFile(logDir, "app-trace.log")
}
