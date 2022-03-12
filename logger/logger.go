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
	"github.com/isyscore/isc-gobase/cron"
	"github.com/isyscore/isc-gobase/time"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	t0 "time"
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

//callerMarshalFunc if you call the Info or Warn etd,the caller will lose it's original caller info,so it will to get it's original caller
//suggest: please use zerolog's Func,such as log.Info,log.Debug and so on,eg:
//log.Info().Msg("%s am a little Cutie","酷达舒")
//log.Debug().Msg("%s say me too","kucs")
func callerMarshalFunc(file string, l int) string {
	if strings.Contains(file, "logger/logger.go") {
		_, f, line, _ := runtime.Caller(6)
		file = f
		l = line
	}
	return file + ":" + strconv.Itoa(l)
}

//InitLog create a root logger. it will write to console and multiple file by level.
// note: default set root logger level is info
// it provides custom log with CustomizeFiles,if it match any caller's name ,log's level will be setting debug and output
func InitLog(logLevel string, timeFmt string, colored bool, appName string) {
	//日志级别设置，默认Info
	zerolog.ErrorHandler = func(err error) {
		// do nothing
	}

	SetGlobalLevel(logLevel)

	initLogDir()
	zerolog.CallerSkipFrameCount = 2
	zerolog.CallerMarshalFunc = callerMarshalFunc
	//时间格式设置
	zerolog.TimeFieldFormat = timeFmt
	//设置日志输出
	out := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: colored}
	out.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf(" [%s] [%-2s]", appName, i))
	}

	outers := append(listWriter, out)
	writer := zerolog.MultiLevelWriter(outers...)

	log.Logger = log.Logger.Output(writer).With().Caller().Timestamp().Logger()

}

type FileLevelWriter struct {
	*os.File
	level zerolog.Level
}

func (lw *FileLevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level.String() == lw.level.String() {
		return lw.Write(p)
	}
	return 0, nil
}

var listWriter []io.Writer

//initLoggerFile open or create and open log file
func initLoggerFile(logDir string, fileName string, level zerolog.Level) {
	pwd, _ := os.Getwd()
	linkName := filepath.Join(pwd, logDir, fileName)
	var file *os.File
	fileHandler := func() {
		t := t0.Now()
		//关闭现有的流
		for _, w := range listWriter {
			if w == nil {
				continue
			}
			if fw, ok := w.(FileLevelWriter); ok {
				fw.Close()
			}
		}
		d, _ := t0.ParseDuration("-24h")
		t1 := t.Add(d)
		//时间格式转换
		logFile := strings.ReplaceAll(linkName, ".log", fmt.Sprintf("-%s.log", time.TimeToStringFormat(t1, time.FmtYMdHms)))
		if err := os.Rename(linkName, logFile); err != nil {
			fmt.Printf("%v\n", err)
			log.Printf("logfile[%s]重命名失败%v", logFile, err)
		}
		//打开创建流
		file, _ = os.OpenFile(linkName, os.O_CREATE, 0666)
		go func() {
			//判断文件是否为空，如果为空则删除
			fi, _ := os.Stat(logFile)
			if fi != nil {
				println("logFile[%s]文件大小:%v", fi.Size())
				if fi.Size() == 0 {
					os.Remove(logFile)
				}
			}
		}()
	}
	fileHandler()
	//每天创建一个文件
	c_d := cron.New()
	c_d.AddFunc("*/5 * * * * ?", fileHandler)
	c_d.Start()
	if file != nil {
		fw := &FileLevelWriter{file, level}
		listWriter = append(listWriter, fw)
	}
}

//initLogDir create log dir and file
func initLogDir() {
	// 创建日志目录
	logDir := filepath.Join(".", "logs")
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.Mkdir(logDir, os.ModePerm)
	}
	// 创建日志文件
	initLoggerFile(logDir, "app-info.log", zerolog.InfoLevel)
	initLoggerFile(logDir, "app-debug.log", zerolog.DebugLevel)
	initLoggerFile(logDir, "app-warn.log", zerolog.WarnLevel)
	initLoggerFile(logDir, "app-error.log", zerolog.ErrorLevel)
	initLoggerFile(logDir, "app-assert.log", zerolog.NoLevel)
	initLoggerFile(logDir, "app-trace.log", zerolog.TraceLevel)
}
