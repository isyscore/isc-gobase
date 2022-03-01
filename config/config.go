package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var appProperty *ApplicationProperty

func LoadConfig() {
	LoadConfigFromRelativePath("")
}

// LoadConfigFromRelativePath 加载相对文件路径
func LoadConfigFromRelativePath(resourceAbsPath string) {
	dir, _ := os.Getwd()
	pkg := strings.Replace(dir, "\\", "/", -1)

	LoadConfigFromAbsPath(path.Join(pkg, "", resourceAbsPath))
}

// LoadConfigFromAbsPath 加载绝对文件路径
func LoadConfigFromAbsPath(resourceAbsPath string) {
	doLoadConfigFromAbsPath(resourceAbsPath)

	// 读取cm文件
	//AppendConfigFromAbsPath("/home/" + GetValueString("base.application.name") + "/config/application-default.yml")
	AppendConfigFromRelativePath("./config/application-default.yml")

	// 加载内部配置
	if err := GetValueObject("server", &ServerCfg); err != nil {
		log.Error().Msgf("加载 Server 配置失败(%v)", err)
	}
	if err := GetValueObject("base", &BaseCfg); err != nil {
		log.Error().Msgf("加载 Base 配置失败(%v)", err)
	}

	if err := GetValueObject("log", &LogCfg); err != nil {
		log.Error().Msgf("加载 Log 配置失败(%v)", err)
	}
}

// AppendConfigFromRelativePath 追加配置：相对路径的配置文件
func AppendConfigFromRelativePath(fileName string) {
	dir, _ := os.Getwd()
	pkg := strings.Replace(dir, "\\", "/", -1)
	fileName = path.Join(pkg, "", fileName)
	extend := getFileExtension(fileName)
	extend = strings.ToLower(extend)
	switch extend {
	case "yaml":
		AppendYamlFile(fileName)
	case "yml":
		AppendYamlFile(fileName)
	case "properties":
		AppendPropertyFile(fileName)
	case "json":
		AppendJsonFile(fileName)
	}
}

// AppendConfigFromAbsPath 追加配置：绝对路径的配置文件
func AppendConfigFromAbsPath(fileName string) {
	extend := getFileExtension(fileName)
	extend = strings.ToLower(extend)
	switch extend {
	case "yaml":
		AppendYamlFile(fileName)
	case "yml":
		AppendYamlFile(fileName)
	case "properties":
		AppendPropertyFile(fileName)
	case "json":
		AppendJsonFile(fileName)
	}
}

type EnvProperty struct {
	Key   string
	Value string
}

func GetConfigValues(c *gin.Context) {
	if nil != appProperty {
		c.Data(200, "application/json; charset=utf-8", []byte(isc.ObjectToJson(appProperty.ValueMap)))
	} else {
		c.Data(200, "application/json; charset=utf-8", []byte("{}"))
	}
}

func GetConfigValue(c *gin.Context) {
	if nil != appProperty {
		value := GetValue(c.Param("key"))
		if nil == value {
			c.Data(200, "application/json; charset=utf-8", []byte(""))
			return
		}
		if isc.IsBaseType(reflect.TypeOf(value)) {
			c.Data(200, "application/json; charset=utf-8", []byte(isc.ToString(value)))
		} else {
			c.Data(200, "application/json; charset=utf-8", []byte(isc.ObjectToJson(value)))
		}
	} else {
		c.Data(200, "application/json; charset=utf-8", []byte("{}"))
	}
}

func UpdateConfig(c *gin.Context) {
	envProperty := EnvProperty{}
	err := isc.DataToObject(c.Request.Body, &envProperty)
	if err != nil {
		log.Warn().Msgf("解析失败，%v", err.Error())
		return
	}

	SetValue(envProperty.Key, envProperty.Value)
}

// 多种格式优先级：json > properties > yaml > yml
func doLoadConfigFromAbsPath(resourceAbsPath string) {
	if !strings.HasSuffix(resourceAbsPath, "/") {
		resourceAbsPath += "/"
	}
	files, err := ioutil.ReadDir(resourceAbsPath)
	if err != nil {
		log.Warn().Msgf("读取配置资源失败，路径(%v), 异常(%v)", resourceAbsPath, err.Error())
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	profile := getActiveProfile()

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if !strings.HasPrefix(fileName, "application") {
			continue
		}

		// 默认配置
		if "application.yaml" == fileName {
			LoadYamlFile(resourceAbsPath + "application.yaml")
			break
		} else if "application.yml" == fileName {
			LoadYamlFile(resourceAbsPath + "application.yml")
			break
		} else if "application.properties" == fileName {
			LoadPropertyFile(resourceAbsPath + "application.properties")
			break
		} else if "application.json" == fileName {
			LoadJsonFile(resourceAbsPath + "application.json")
			break
		}

		if "" != profile {
			currentProfile := getProfileFromFileName(fileName)
			if currentProfile == profile {
				extend := getFileExtension(fileName)
				extend = strings.ToLower(extend)
				if "yaml" == extend {
					LoadYamlFile(resourceAbsPath + fileName)
					break
				} else if "yml" == extend {
					LoadYamlFile(resourceAbsPath + fileName)
					break
				} else if "properties" == extend {
					LoadPropertyFile(resourceAbsPath + fileName)
					break
				} else if "json" == extend {
					LoadJsonFile(resourceAbsPath + fileName)
					break
				}
			}
		}
	}
	SetValue("base.actives.profile", profile)
}

// 临时写死
// 优先级：本地配置 > 启动参数 > 环境变量
func getActiveProfile() string {
	profile := GetValueString("base.actives.profile")
	if "" != profile {
		return profile
	}

	flag.StringVar(&profile, "base.actives.profile", "", "环境变量")
	flag.Parse()
	if "" != profile {
		return profile
	}

	profile = os.Getenv("base.actives.profile")
	if "" != profile {
		return profile
	}
	return ""
}

func getProfileFromFileName(fileName string) string {
	if strings.HasPrefix(fileName, "application-") {
		words := strings.SplitN(fileName, ".", 2)
		appNames := words[0]

		appNameAndProfile := strings.SplitN(appNames, "-", 2)
		return appNameAndProfile[1]
	}
	return ""
}

func getFileExtension(fileName string) string {
	if strings.Contains(fileName, ".") {
		words := strings.SplitN(fileName, ".", 2)
		return words[1]
	}
	return ""
}

func LoadYamlFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	property, err := isc.YamlToProperties(string(content))
	valueMap, _ := isc.PropertiesToMap(property)
	appProperty.ValueMap = valueMap

	yamlMap, err := isc.YamlToMap(string(content))
	appProperty.ValueDeepMap = yamlMap
}

func AppendYamlFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	property, err := isc.YamlToProperties(string(content))
	valueMap, _ := isc.PropertiesToMap(property)
	for k, v := range valueMap {
		SetValue(k, v)
	}
}

func LoadPropertyFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	valueMap, _ := isc.PropertiesToMap(string(content))
	appProperty.ValueMap = valueMap

	yamlStr, _ := isc.PropertiesToYaml(string(content))
	yamlMap, _ := isc.YamlToMap(yamlStr)
	appProperty.ValueDeepMap = yamlMap
}

func AppendPropertyFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	valueMap, _ := isc.PropertiesToMap(string(content))
	for k, v := range valueMap {
		SetValue(k, v)
	}
}

func LoadJsonFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	yamlStr, _ := isc.JsonToYaml(string(content))
	property, _ := isc.YamlToProperties(yamlStr)
	valueMap, _ := isc.PropertiesToMap(property)
	appProperty.ValueMap = valueMap

	yamlMap, _ := isc.YamlToMap(yamlStr)
	appProperty.ValueDeepMap = yamlMap
}

func AppendJsonFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn().Msgf("fail to read file:", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	yamlStr, err := isc.JsonToYaml(string(content))
	property, err := isc.YamlToProperties(yamlStr)
	valueMap, _ := isc.PropertiesToMap(property)
	for k, v := range valueMap {
		SetValue(k, v)
	}
}

func SetValue(key string, value interface{}) {
	if nil == value {
		return
	}
	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = map[string]interface{}{}
	}

	if oldValue, exist := appProperty.ValueMap[key]; exist {
		if reflect.TypeOf(oldValue) != reflect.TypeOf(oldValue) {
			return
		}
		appProperty.ValueMap[key] = value
	}
	doPutValue(key, value)
}

func doPutValue(key string, value interface{}) {
	if strings.Contains(key, ".") {
		oldValue := GetValue(key)
		if nil == oldValue {
			return
		}
		if reflect.TypeOf(oldValue).Kind() != reflect.TypeOf(value).Kind() {
			return
		}

		lastIndex := strings.LastIndex(key, ".")
		startKey := key[:lastIndex]
		endKey := key[lastIndex+1:]

		data := GetValue(startKey)
		startValue := isc.ToMap(data)
		if nil != startValue {
			startValue[endKey] = value
		}

		doPutValue(startKey, startValue)
	}
	appProperty.ValueDeepMap[key] = value
}

func GetValueString(key string) string {
	if nil == appProperty {
		return ""
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToString(value)
	}
	return ""
}

func GetValueInt(key string) int {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt(value)
	}
	return 0
}

func GetValueInt8(key string) int8 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt8(value)
	}
	return 0
}

func GetValueInt16(key string) int16 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt16(value)
	}
	return 0
}

func GetValueInt32(key string) int32 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt32(value)
	}
	return 0
}

func GetValueInt64(key string) int64 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt64(value)
	}
	return 0
}

func GetValueUInt(key string) uint {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt(value)
	}
	return 0
}

func GetValueUInt8(key string) uint8 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt8(value)
	}
	return 0
}

func GetValueUInt16(key string) uint16 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt16(value)
	}
	return 0
}

func GetValueUInt32(key string) uint32 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt32(value)
	}
	return 0
}

func GetValueUInt64(key string) uint64 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt64(value)
	}
	return 0
}

func GetValueFloat32(key string) float32 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToFloat32(value)
	}
	return 0
}

func GetValueFloat64(key string) float64 {
	if nil == appProperty {
		return 0
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToFloat64(value)
	}
	return 0
}

func GetValueBool(key string) bool {
	if nil == appProperty {
		return false
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToBool(value)
	}
	return false
}

func GetValueStringDefault(key, defaultValue string) string {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToString(value)
	}
	return defaultValue
}

func GetValueIntDefault(key string, defaultValue int) int {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt(value)
	}
	return defaultValue
}

func GetValueInt8Default(key string, defaultValue int8) int8 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt8(value)
	}
	return defaultValue
}

func GetValueInt16Default(key string, defaultValue int16) int16 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt16(value)
	}
	return defaultValue
}

func GetValueInt32Default(key string, defaultValue int32) int32 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt32(value)
	}
	return defaultValue
}

func GetValueInt64Default(key string, defaultValue int64) int64 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToInt64(value)
	}
	return defaultValue
}

func GetValueUIntDefault(key string, defaultValue uint) uint {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt(value)
	}
	return defaultValue
}

func GetValueUInt8Default(key string, defaultValue uint8) uint8 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt8(value)
	}
	return defaultValue
}

func GetValueUInt16Default(key string, defaultValue uint16) uint16 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt16(value)
	}
	return defaultValue
}

func GetValueUInt32Default(key string, defaultValue uint32) uint32 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt32(value)
	}
	return defaultValue
}

func GetValueUInt64Default(key string, defaultValue uint64) uint64 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToUInt64(value)
	}
	return defaultValue
}

func GetValueFloat32Default(key string, defaultValue float32) float32 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToFloat32(value)
	}
	return defaultValue
}

func GetValueFloat64Default(key string, defaultValue float64) float64 {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToFloat64(value)
	}
	return defaultValue
}

func GetValueBoolDefault(key string, defaultValue bool) bool {
	if nil == appProperty {
		return defaultValue
	}
	if value, exist := appProperty.ValueMap[key]; exist {
		return isc.ToBool(value)
	}
	return defaultValue
}

func GetValueObject(key string, targetPtrObj interface{}) error {
	if nil == appProperty {
		return nil
	}
	data := doGetValue(appProperty.ValueDeepMap, key)
	err := isc.DataToObject(data, targetPtrObj)
	if err != nil {
		return err
	}
	return nil
}

func GetValue(key string) interface{} {
	if nil == appProperty {
		return nil
	}
	return doGetValue(appProperty.ValueDeepMap, key)
}

func doGetValue(parentValue interface{}, key string) interface{} {
	if key == "" {
		return parentValue
	}
	parentValueKind := reflect.ValueOf(parentValue).Kind()
	if parentValueKind == reflect.Map {
		keys := strings.SplitN(key, ".", 2)
		v1 := reflect.ValueOf(parentValue).MapIndex(reflect.ValueOf(keys[0]))
		emptyValue := reflect.Value{}
		if v1 == emptyValue {
			return nil
		}
		if len(keys) == 1 {
			return doGetValue(v1.Interface(), "")
		} else {
			return doGetValue(v1.Interface(), fmt.Sprintf("%v", keys[1]))
		}
	}
	return nil
}

type ApplicationProperty struct {
	ValueMap     map[string]interface{}
	ValueDeepMap map[string]interface{}
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
		log.Error().Msgf("读取文件异常(%v)", err)
	}
	return handler(data, AConfig)
}

//LoadSpringConfig read fileName from current dictionary and fileName is application.yml,eg:/home/isc-gobase/application.yml, and transform it to AConfig
//note: AConfig must be a pointer
//note: if it has Spring.Profiles.Active,eg: Spring.Profiles.Active=dev,will load config from /home/isc-gobase/application-dev.yml,and same key
//will write in the last one.

func LoadSpringConfig(AConfig any) {
	_ = LoadYamlConfig("application.yml", AConfig, func(data []byte, AConfig any) error {
		err := yaml.Unmarshal(data, AConfig)
		if err != nil {
			log.Error().Msgf("读取 application.yml 异常(%v)", err)
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
				log.Error().Msgf("读取 application-%s.yml 失败", act)
				return err
			} else {
				err = yaml.Unmarshal(yamlAdditional, AConfig)
				if err != nil {
					log.Error().Msgf("读取 application-%s.yml 异常", act)
					return err
				}
			}
		}
		return nil
	})
}
