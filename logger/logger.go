package logger

import (
	"log"
	"os"
	"path/filepath"
)

func doLog(l *log.Logger, format string, v ...any) {
	log.Printf(format+"\n", v...)
	if l != nil {
		l.Printf(format+"\n", v...)
	}
}

func Info(format string, v ...any) {
	doLog(loggerInfo, format, v...)
}

func Warn(format string, v ...any) {
	doLog(loggerWarn, format, v...)
}

func Error(format string, v ...any) {
	doLog(loggerError, format, v...)
}

func Debug(format string, v ...any) {
	doLog(loggerDebug, format, v...)
}

func Assert(format string, v ...any) {
	doLog(loggerAssert, format, v...)
}

var loggerInfo *log.Logger = nil
var loggerDebug *log.Logger = nil
var loggerWarn *log.Logger = nil
var loggerError *log.Logger = nil
var loggerAssert *log.Logger = nil

func initLoggerFile(logDir string, fileName string) *log.Logger {
	var l *log.Logger
	logFile := filepath.Join(logDir, fileName)
	if file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err == nil {
		l = log.New(file, "", log.LstdFlags|log.Llongfile)
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
