package test

import (
	"fmt"
	"testing"

	"github.com/isyscore/isc-gobase/config"
)

func TestRead(t *testing.T) {
	config.LoadConfig()
	fmt.Println(config.GetValueString("pvcStorage"))
}
