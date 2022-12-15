package test

import (
	"github.com/isyscore/isc-gobase/baselog"
	"github.com/isyscore/isc-gobase/config"
	"testing"
)

func TestInfo(t *testing.T) {
	config.LoadYamlFile("./application.yaml")
	baselog.InitLog()

	// info
	baselog.Info("hello %v", "info")
	baselog.GetLogger("group1").Info("hello", " ","info")
	baselog.GetLogger("group1").Infof("hello %v", "info")
}

func TestLevel(t *testing.T) {
	// debug
	baselog.Debug("hello %v", "debug")

	// warn
	baselog.Warn("hello %v", "warn")
}
