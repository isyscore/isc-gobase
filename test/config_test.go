package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"os"
	"testing"
)

func TestReadDefault(t *testing.T) {
	//config.LoadConfig()
	config.LoadConfigFromRelativePath("./resources/")

	err := config.GetValueObject("server", &config.BaseCfg)
	if err != nil {
		return
	}

	fmt.Println(isc.ObjectToJson(config.BaseCfg))

	fmt.Println(os.Environ())
}

func TestReadDefault2(t *testing.T) {
	config.LoadConfig()

	fmt.Println(config.GetValueString("one.name"))
	fmt.Println("xxx")
}
