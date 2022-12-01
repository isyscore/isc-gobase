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
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/listener"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	t0 "time"

	"github.com/isyscore/isc-gobase/cron"
	"github.com/isyscore/isc-gobase/time"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var loggerTask *cron.Cron

func init() {
	cfg := config.BaseLogger{}
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Time.Format == "" {
		cfg.Time.Format = "2006-01-02 15:04:05"
	}
	if cfg.Split.Size == 0 {
		cfg.Split.Size = 300
	}
	if cfg.Max.History == 0 {
		cfg.Max.History = 7
	}

	appName := ""

	//日志级别设置，默认Info
	zerolog.ErrorHandler = func(err error) {
		// do nothing
	}

	SetGlobalLevel(cfg.Level)

	zerolog.CallerSkipFrameCount = 2
	zerolog.CallerMarshalFunc = callerMarshalFunc
	zerolog.TimeFieldFormat = cfg.Time.Format
	out := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: cfg.Color.Enable, FormatTimestamp: func(i interface{}) string {
		return "[" + time.Now().Format(cfg.Time.Format) + "]"
	}}
	out.FormatLevel = func(i any) string {
		return fmt.Sprintf("[%s] [%v] [%v] [%-2s]", appName, GetMdc(constants.TRACE_HEAD_ID), GetMdc(constants.TRACE_HEAD_USER_ID), i)
	}
	out.FormatCaller = callerFormatter
	initLogDir(out, cfg.Split.Enable, cfg.Split.Size, cfg.Dir, cfg.Max.History, appName, cfg.Console.WriteFile)
}

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
	log.WithLevel(zerolog.Disabled).Msgf(format, v...)
}

func Panic(format string, v ...any) {
	log.WithLevel(zerolog.PanicLevel).Msgf(format, v...)
}

func Fatal(format string, v ...any) {
	log.WithLevel(zerolog.FatalLevel).Msgf(format, v...)
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

// SetGlobalLevel sets the global override for log level. If this
// values is raised, all Loggers will use at least this value.
//
// To globally disable logs, set zerolog.GlobalLevel to be Disabled.
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
// it provides custom log with CustomizeFiles,if it matches any caller's name ,log's level will be setting debug and output
func InitLog(appName string) {
	cfg := config.BaseCfg.Logger
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Time.Format == "" {
		cfg.Time.Format = "2006-01-02 15:04:05"
	}
	if cfg.Split.Size == 0 {
		cfg.Split.Size = 300
	}
	if cfg.Max.History == 0 {
		cfg.Max.History = 7
	}

	//日志级别设置，默认Info
	zerolog.ErrorHandler = func(err error) {
		// do nothing
	}

	SetGlobalLevel(cfg.Level)

	zerolog.CallerSkipFrameCount = 2
	zerolog.CallerMarshalFunc = callerMarshalFunc
	//时间格式设置
	zerolog.TimeFieldFormat = cfg.Time.Format
	//设置日志输出
	out := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: cfg.Color.Enable, FormatTimestamp: func(i interface{}) string {
		return "[" + time.Now().Format(cfg.Time.Format) + "]"
	}}
	out.FormatLevel = func(i any) string {
		return fmt.Sprintf("[%s] [%v] [%v] [%-2s]", appName, GetMdc(constants.TRACE_HEAD_ID), GetMdc(constants.TRACE_HEAD_USER_ID), i)
	}
	out.FormatCaller = callerFormatter
	initLogDir(out, cfg.Split.Enable, cfg.Split.Size, cfg.Dir, cfg.Max.History, appName, cfg.Console.WriteFile)

	// 添加配置变更事件的监听
	listener.AddListener(listener.EventOfConfigChange, ConfigChangeListener)
}

func ConfigChangeListener(event listener.BaseEvent) {
	ev := event.(listener.ConfigChangeEvent)
	if ev.Key == "base.logger.level" {
		SetGlobalLevel(ev.Value)
	}
}

type FileLevelWriter struct {
	*os.File
	level  zerolog.Level
	writer zerolog.ConsoleWriter
}

func (lw *FileLevelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level.String() == lw.level.String() {
		return lw.writer.Write(p)
	}
	return len(p), nil
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
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatal().Msgf("日志目录创建异常:%v", err)
		}
	}
	return logDir
}

var panicHandler = Strategy{}

func callerFormatter(i interface{}) string {
	loggerPath := config.BaseCfg.Logger.Path
	if loggerPath == "" || loggerPath == "short" {
		// 去除过多的目录层级信息
		str := i.(string)
		strs := strings.Split(str, string(os.PathSeparator))
		ret := strs[len(strs)-1]
		if len(strs) > 1 {
			ret = strs[len(strs)-2] + string(os.PathSeparator) + ret
		}
		return ret
	} else if loggerPath == "full" {
		//去除Jenkins或编译所在主机信息
		str := i.(string)
		strs := strings.Split(str, "@2/project")
		if len(strs) > 1 {
			strs[0] = "../.."
			str = strings.Join(strs, "")
		}

		return str
	} else {
		return i.(string)
	}
}

func createFileLeveWriter(level zerolog.Level, strTime string, idx int, dir, appName string) *FileLevelWriter {
	strL := level.String()
	if level == zerolog.Disabled {
		strL = "console"
	}
	linkName := fmt.Sprintf("app-%s.log", strL)
	linkName = filepath.Join(getLogDir(dir), linkName)
	logFile := strings.ReplaceAll(linkName, ".log", fmt.Sprintf("-%s.log", strTime))
	if idx > 0 {
		logFile = strings.ReplaceAll(logFile, ".log", fmt.Sprintf(".%d.log", idx))
	}
	//建立软链

	if _, err := os.Stat(linkName); err != nil {
		if os.IsExist(err) {
			_ = os.Remove(linkName)
		}
	} else {
		_ = os.Remove(linkName)
	}

	if strings.ToLower(runtime.GOOS) != "windows" {
		_ = os.Symlink(logFile, linkName)
	} else {
		_ = os.Link(logFile, linkName)
	}

	//打开创建流
	file1, _ := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR, 0666)

	fw := &FileLevelWriter{file1, level, zerolog.ConsoleWriter{
		Out:     file1,
		NoColor: false,
		FormatTimestamp: func(i interface{}) string {
			return "[" + time.Now().Format(time.FmtYMdHmsSSS) + "]"
		},
		FormatLevel: func(i any) string {
			return fmt.Sprintf("[%s] [%v] [%v] [%-2s]", appName, GetMdc(constants.TRACE_HEAD_ID), GetMdc(constants.TRACE_HEAD_USER_ID), i)
		},

		FormatCaller: callerFormatter,
	}}
	if level == zerolog.PanicLevel {
		if err := panicHandler.Dup2(fw, os.Stderr); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "system panic log redirect to %s failed:%v", logFile, err)
		}
	}
	return fw
}

var levels = []zerolog.Level{zerolog.DebugLevel, zerolog.TraceLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.ErrorLevel, zerolog.PanicLevel, zerolog.FatalLevel, zerolog.Disabled}

func updateOuters(out zerolog.ConsoleWriter, idx int, ls []zerolog.Level, dir, name string, write2File bool) {
	//关闭现有流
	closeFileLevelWriter(oldWriter)
	//修改listWriter
	var newWriter []io.Writer
	//时间格式转换
	strTime := time.TimeToStringFormat(t0.Now(), time.FmtYMd)
	for _, level := range ls {
		fw := createFileLeveWriter(level, strTime, idx, dir, name)
		if fw != nil {
			newWriter = append(newWriter, fw)
		}
		if level == zerolog.Disabled && write2File {
			os.Stderr = fw.File
			os.Stdout = fw.File
		}
		if level == zerolog.PanicLevel {
			if err := panicHandler.Dup2(fw, os.Stderr); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "system panic log redirect to %s file failed:%v", fw.level.String(), err)
			}
		}
	}

	outers := append(newWriter, out)
	writer := zerolog.MultiLevelWriter(outers...)
	log.Logger = log.Logger.Output(writer).With().Caller().Logger()
	oldWriter = newWriter
}

var oldWriter []io.Writer

//initLogDir create log dir and file
func initLogDir(out zerolog.ConsoleWriter, splitEnable bool, splitSize int64, dir string, history int, name string, write2File bool) {
	if loggerTask != nil {
		loggerTask.Stop()
	}
	loggerTask = cron.New()

	// 每天创建一个日志文件
	fileHandler := func() { updateOuters(out, 0, levels, dir, name, write2File) }
	fileHandler()
	_ = loggerTask.AddFunc("0 0 0 * * ?", fileHandler)

	if splitEnable {
		_ = loggerTask.AddFunc("*/1 * * * * ?", func() {
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
							updateOuters(out, idx+1, []zerolog.Level{fw.level}, dir, name, write2File)
						}
					}
				}
			}()
		})
	}

	// log.Info().Msgf("开启定时日志清理任务")
	_ = loggerTask.AddFunc("0 0 1 * * ?", func() {
		log.Debug().Msg("定时每天日志清理任务执行")
		_ = filepath.Walk(getLogDir(dir), func(path string, info fs.FileInfo, err error) error {
			now := time.Now()
			if time.DaysBetween(now, info.ModTime()) > history {
				//remove file
				err = os.Remove(path)
				defer func() error {
					if x := recover(); x != nil {
						log.Error().Msgf("日志文件[%s]删除错误", err)
						return x.(error)
					}
					return nil
				}()
			}
			return err
		})
	})

	loggerTask.Start()
}
