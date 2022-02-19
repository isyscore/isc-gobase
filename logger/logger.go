package logger

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

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
	out.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" [%s] [%-2s]", config.GetValueStringDefault("base.application.name", "isc-gobase"), i))
	}
	log.Logger = log.Logger.Output(out).With().Caller().Logger()
	//添加hook
	//
	//levelInfoHook := zerolog.HookFunc(func(e *zerolog.Event, l zerolog.Level, msg string){
	//	levelName := l.String()
	//	if l == zerolog.NoLevel {
	//		e.Discard()
	//	} else if l == zerolog.InfoLevel {
	//		log.Logger.Output()
	//	}
	//}
}

//func initLoggerFile(logDir string, fileName string) *log.Logger {
//	var l *log.Logger
//	logFile := filepath.Join(logDir, fileName)
//	if file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err == nil {
//		l = log.New(file, "", log.LstdFlags|log.Llongfile)
//	}
//	return l
//}
//
//func init() {
//	// 创建日志目录
//	logDir := filepath.Join(".", "logs")
//	if _, err := os.Stat(logDir); os.IsNotExist(err) {
//		_ = os.Mkdir(logDir, os.ModePerm)
//	}
//	// 创建日志文件
//	loggerInfo = initLoggerFile(logDir, "app-info.log")
//	loggerDebug = initLoggerFile(logDir, "app-debug.log")
//	loggerWarn = initLoggerFile(logDir, "app-warn.log")
//	loggerError = initLoggerFile(logDir, "app-error.log")
//	loggerAssert = initLoggerFile(logDir, "app-assert.log")
//
//}
