package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"testing"
)

func TestReadDefault(t *testing.T) {
	//config.LoadConfigFromRelativePath("./resources/")

	//err := config.GetValueObject("base", &config.BaseCfg)
	//if err != nil {
	//	return
	//}

	fmt.Println(config.ApiModule)
	fmt.Println(config.BaseCfg)
}

func TestReadDefault2(t *testing.T) {
	config.LoadConfig()

	//fmt.Println(config.GetValueString("one.name"))
	fmt.Println(config.GetValueArrayInt("base.server.exception.print.except"))
}
