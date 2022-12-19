package test

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestInfo(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	logger.InitLog()

	// info
	logger.Info("hello %v", "info")
	logger.Group("group1").Info("hello", " ", "info")
	logger.Group("group1").Infof("hello %v", "info")
}

func TestLevel(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	logger.InitLog()

	// debug
	logger.Debug("hello %v", "debug")

	// info
	logger.Info("hello %v", "info")

	// warn
	logger.Warn("hello %v", "warn")

	// error
	logger.Error("hello %v", "error")
}

func TestLevelChange(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	logger.InitLog()

	// info
	logger.Info("hello %v", "info1")
	logger.Info("hello %v", "info2")

	// 设置后下面的不再显示
	logger.SetGlobalLevel("warn")
	logger.Info("hello %v", "info3")
}

// 日志分组的级别变更
func TestGroupLevelChange1(t *testing.T) {
	config.LoadYamlFile("./application-debug.yaml")
	logger.InitLog()

	// info
	logger.Info("hello %v", "info1")
	logger.Group("group1").Infof("hello %v", "group1 info1")
	logger.Group("group2").Infof("hello %v", "group2 info2")

	// 设置后下面的不再显示
	logger.SetGlobalLevel("warn")
	// 不再打印
	logger.Info("hello %v", "info3")
	// 继续打印
	logger.Group("group1").Infof("hello %v", "group1 info1")
	logger.Group("group2").Infof("hello %v", "group2 info2")

	warnLevel, _ := logrus.ParseLevel("warn")
	logger.Group("group1").SetLevel(warnLevel)
	// 不再打印
	logger.Group("group1").Infof("hello %v", "group1 change info1")
	// 继续打印
	logger.Group("group2").Infof("hello %v", "group2 change info2")
}

func TestLoggerPathShort(t *testing.T) {
	config.LoadYamlFile("./application-short.yaml")
	logger.InitLog()

	logger.Info("test")
}

func TestLoggerPathFull(t *testing.T) {
	config.LoadYamlFile("./application-full.yaml")
	logger.InitLog()

	logger.Info("test")
}

func TestLoggerRotate(t *testing.T) {
	config.LoadYamlFile("./application-rotate.yaml")
	logger.InitLog()

	//for i := 0; i < 100; i++ {
	//	logger.Info("test")
	//	time.Sleep(1 * time.Second)
	//}
}
