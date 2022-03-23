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
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	t0 "time"

	"github.com/isyscore/isc-gobase/cron"
	f0 "github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/time"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	log.WithLevel(zerolog.NoLevel).Msgf(format, v...)
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
func InitLog(logLevel string, timeFmt string, colored bool, appName string, splitEnable bool, splitSize int64, logDir string) {
	//日志级别设置，默认Info
	zerolog.ErrorHandler = func(err error) {
		// do nothing
	}

	zerolog.CallerMarshalFunc = callerMarshalFunc

	SetGlobalLevel(logLevel)

	zerolog.CallerSkipFrameCount = 2
	zerolog.CallerMarshalFunc = callerMarshalFunc
	//时间格式设置
	zerolog.TimeFieldFormat = timeFmt
	//设置日志输出
	out := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: colored}
	out.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf(" [%s] [%-2s]", appName, i))
	}
	initLogDir(out, splitEnable, splitSize, logDir)

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

func closeFileLevelWriter(writers []io.Writer) {
	if len(writers) == 0 {
		return
	}
	for _, w := range writers {
		if w == nil {
			continue
		}
		if fw, ok := w.(*FileLevelWriter); ok {
			_ = fw.Close()
			fi, _ := os.Stat(fw.File.Name())
			if fi != nil && fi.Size() == 0 {
				_ = os.Remove(fw.File.Name())
			}
		}
	}
}

func getLogDir(logDir string) string {
	if logDir == "" || !strings.HasPrefix(logDir, "/") {
		pwd, _ := os.Getwd()
		// 创建日志目录
		logDir = filepath.Join(pwd, "logs")
	}

	fmt.Println("日志目录", logDir)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.Mkdir(logDir, os.ModePerm)
	}
	return logDir
}

func createFileLeveWriter(level zerolog.Level, strTime string, idx int, dir string) *FileLevelWriter {
	strL := level.String()
	if level == zerolog.NoLevel {
		strL = "assert"
	}
	linkName := fmt.Sprintf("app-%s.log", strL)
	linkName = filepath.Join(getLogDir(dir), linkName)
	logFile := strings.ReplaceAll(linkName, ".log", fmt.Sprintf("-%s.log", strTime))
	if idx > 0 {
		logFile = strings.ReplaceAll(logFile, ".log", fmt.Sprintf(".%d.log", idx))
	}
	//建立软链

	if f0.FileExists(linkName) {
		f0.DeleteFile(linkName)
	}

	if strings.ToLower(runtime.GOOS) != "windows" {
		_ = os.Symlink(logFile, linkName)
	} else {
		_ = os.Link(logFile, linkName)
	}

	//打开创建流
	file1, _ := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR, 0644)
	return &FileLevelWriter{file1, level}

}

var levels = []zerolog.Level{zerolog.NoLevel, zerolog.DebugLevel, zerolog.TraceLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.FatalLevel, zerolog.PanicLevel}

func updateOuters(out zerolog.ConsoleWriter, idx int, ls []zerolog.Level, dir string) {
	//关闭现有流
	closeFileLevelWriter(oldWriter)
	//修改listWriter
	var newWriter []io.Writer
	//时间格式转换
	strTime := time.TimeToStringFormat(t0.Now(), time.FmtYMd)
	for _, level := range ls {
		fw := createFileLeveWriter(level, strTime, idx, dir)
		newWriter = append(newWriter, fw)
	}
	outers := append(newWriter, out)
	writer := zerolog.MultiLevelWriter(outers...)
	log.Logger = log.Logger.Output(writer).With().Caller().Logger()
	oldWriter = newWriter
}

var oldWriter []io.Writer

//initLogDir create log dir and file
func initLogDir(out zerolog.ConsoleWriter, splitEnable bool, splitSize int64, dir string) {
	fileHandler := func() { updateOuters(out, 0, levels, dir) }
	fileHandler()

	//每天创建一个文件
	c_d := cron.New()
	_ = c_d.AddFunc("0 0 0 * * ?", fileHandler)
	c_d.Start()

	if splitEnable {
		c_check := cron.New()
		_ = c_check.AddFunc("*/1 * * * * ?", func() {
			//检查文件大小，如果超过核定大小，则生成新文件
			go func() {
				for _, w := range oldWriter {
					if fw, ok := w.(*FileLevelWriter); ok {
						//判断文件大小,默认300M
						if fi, _ := os.Stat(fw.File.Name()); fi != nil && fi.Size() >= (splitSize<<20) {
							name := fi.Name()
							idxs := strings.Split(name, ".")
							idx := 0
							if len(idxs) == 3 {
								idx, _ = strconv.Atoi(idxs[1])
							} else if len(idxs) > 3 {
								idx, _ = strconv.Atoi(idxs[len(idxs)-2])
							}
							updateOuters(out, idx+1, []zerolog.Level{fw.level}, dir)
						}
					}
				}
			}()
		})
		c_check.Start()
	}
}
