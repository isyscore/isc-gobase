package test

import (
	"github.com/isyscore/isc-gobase/baselog"
	"github.com/isyscore/isc-gobase/config"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestInfo(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	baselog.InitLog()

	// info
	baselog.Info("hello %v", "info")
	baselog.GetLogger("group1").Info("hello", " ", "info")
	baselog.GetLogger("group1").Infof("hello %v", "info")
}

func TestLevel(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	baselog.InitLog()

	// debug
	baselog.Debug("hello %v", "debug")

	// info
	baselog.Info("hello %v", "info")

	// warn
	baselog.Warn("hello %v", "warn")

	// error
	baselog.Error("hello %v", "error")
}

func TestLevelChange(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	baselog.InitLog()

	// info
	baselog.Info("hello %v", "info1")
	baselog.Info("hello %v", "info2")

	// 设置后下面的不再显示
	baselog.SetGlobalLevel("warn")
	baselog.Info("hello %v", "info3")
}

// 日志分组的级别变更
func TestGroupLevelChange1(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	baselog.InitLog()

	// info
	baselog.Info("hello %v", "info1")
	baselog.GetLogger("group1").Infof("hello %v", "group1 info1")
	baselog.GetLogger("group2").Infof("hello %v", "group2 info2")

	// 设置后下面的不再显示
	baselog.SetGlobalLevel("warn")
	// 不再打印
	baselog.Info("hello %v", "info3")
	// 继续打印
	baselog.GetLogger("group1").Infof("hello %v", "group1 info1")
	baselog.GetLogger("group2").Infof("hello %v", "group2 info2")

	warnLevel, _ := logrus.ParseLevel("warn")
	baselog.GetLogger("group1").SetLevel(warnLevel)
	// 不再打印
	baselog.GetLogger("group1").Infof("hello %v", "group1 change info1")
	// 继续打印
	baselog.GetLogger("group2").Infof("hello %v", "group2 change info2")
}

func TestLoggerPathShort(t *testing.T) {
	config.LoadYamlFile("./application-short.yaml")
	baselog.InitLog()

	baselog.Info("test")
}

func TestLoggerPathFull(t *testing.T) {
	config.LoadYamlFile("./application-full.yaml")
	baselog.InitLog()

	baselog.Info("test")
}

func TestLoggerRotate(t *testing.T) {
	config.LoadYamlFile("./application-rotate.yaml")
	baselog.InitLog()

	//for i := 0; i < 100; i++ {
	//	baselog.Info("test")
	//	time.Sleep(1 * time.Second)
	//}
}
