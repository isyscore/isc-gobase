package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"isc-gobase/logger"
	"reflect"
)

type AppServer struct {
	Port   int  `yaml:"port"`
	Lookup bool `yaml:"lookup"`
}

type AppSpring struct {
	Application AppApplication `yaml:"application"`
	Profiles    AppProfile     `yaml:"profiles"`
}

type AppLogger struct {
	Level string `yaml:"level"`
}

type AppProfile struct {
	Active string `yaml:"active"`
}

type AppApplication struct {
	Name string `yaml:"name"`
}

func LoadConfig(AConfig any) {
	yamlFile, err := ioutil.ReadFile("./application.yml")
	if err != nil {
		logger.Error("读取 application.yml 失败")
		return
	}
	err = yaml.Unmarshal(yamlFile, AConfig)
	if err != nil {
		logger.Error("读取 application.yml 异常(%v)", err)
		return
	}
	v1 := reflect.ValueOf(AConfig).Elem()
	o1 := v1.FieldByName("Spring").Interface()
	v2 := reflect.ValueOf(o1)
	o2 := v2.FieldByName("Profiles").Interface()
	v3 := reflect.ValueOf(o2)
	act := v3.FieldByName("Active").String()

	if act != "" && act != "default" {
		yamlAdditional, err := ioutil.ReadFile(fmt.Sprintf("./application-%s.yml", act))
		if err != nil {
			logger.Error("读取 application-%s.yml 失败", act)
		} else {
			err = yaml.Unmarshal(yamlAdditional, AConfig)
			if err != nil {
				logger.Error("读取 application-%s.yml 异常", act)
			}
		}
	}
}
