package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/isyscore/isc-gobase/logger"
	"gopkg.in/yaml.v2"
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

//LoadYamlConfig read fileName from private path fileName,eg:application.yml, and transform it to AConfig
//note: AConfig must be a pointer
func LoadYamlConfig(fileName string, AConfig any, handler func(data []byte, AConfig any) error) error {
	pwd, _ := os.Getwd()
	fp := filepath.Join(pwd, fileName)
	return LoadYamlConfigByAbsolutPath(fp, AConfig, handler)
}

//LoadYamlConfigByAbsolutPath read fileName from absolute path fileName,eg:/home/isc-gobase/application.yml, and transform it to AConfig
//note: AConfig must be a pointer
func LoadYamlConfigByAbsolutPath(path string, AConfig any, handler func(data []byte, AConfig any) error) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().Msgf("文件[%s]读取异常\n%v")
	}
	return handler(data, AConfig)
}

//LoadConfig read fileName from current dictionary and fileName is application.yml,eg:/home/isc-gobase/application.yml, and transform it to AConfig
//note: AConfig must be a pointer
//note: if it has Spring.Profiles.Active,eg: Spring.Profiles.Active=dev,will load config from /home/isc-gobase/application-dev.yml,and same key
//will write in the last one.
func LoadConfig(AConfig any) {
	LoadYamlConfig("application.yml", AConfig, func(data []byte, AConfig any) error {
		err := yaml.Unmarshal(data, AConfig)
		if err != nil {
			logger.Error("读取 application.yml 异常(%v)", err)
			return err
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
				return err
			} else {
				err = yaml.Unmarshal(yamlAdditional, AConfig)
				if err != nil {
					logger.Error("读取 application-%s.yml 异常", act)
					return err
				}
			}
		}
		return nil
	})
}
