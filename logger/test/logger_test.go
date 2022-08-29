package test

import (
	"github.com/isyscore/isc-gobase/logger"
	"testing"
)

func init() {
	logger.Info("ttt")
}

func TestInfo(t *testing.T) {
	logger.Info("ok")
}
