package config

import (
	"fmt"
	"github.com/isyscore/isc-gobase/listener"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/file"

	"github.com/isyscore/isc-gobase/isc"
	"gopkg.in/yaml.v2"
)

var appProperty *ApplicationProperty
var configExist = false
var loadLock sync.Mutex
var configLoaded = false
var CurrentProfile = ""

func LoadConfig() {
	loadLock.Lock()
	defer loadLock.Unlock()
	if configLoaded {
		return
	}

	LoadConfigFromRelativePath("")
	configLoaded = true
}

func LoadConfigFromRelativePath(resourceAbsPath string) {
	dir, _ := os.Getwd()
	pkg := strings.Replace(dir, "\\", "/", -1)

	LoadConfigFromAbsPath(path.Join(pkg, "", resourceAbsPath))
}

func LoadConfigFromAbsPath(resourceAbsPath string) {
	doLoadConfigFromAbsPath(resourceAbsPath)

	cmPath := os.Getenv("base.config.additional-location")
	if cmPath == "" {
		cmPath = "./config/application-default.yml"
	}
	AppendConfigFromRelativePath(cmPath)

	ApiModule = GetValueString("api-module")

	if err := GetValueObject("base", &BaseCfg); err != nil {
		log.Printf("加载 Base 配置失败(%v)", err)
	}
}

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

func ExistConfigFile() bool {
	return configExist
}

func GetConfigValues(c *gin.Context) {
	if nil != appProperty {
		c.Data(200, "application/json; charset=utf-8", []byte(isc.ObjectToJson(appProperty.ValueMap)))
	} else {
		c.Data(200, "application/json; charset=utf-8", []byte("{}"))
	}
}

func GetConfigDeepValues(c *gin.Context) {
	if nil != appProperty {
		c.Data(200, "application/json; charset=utf-8", []byte(isc.ObjectToJson(appProperty.ValueDeepMap)))
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
		log.Printf("解析失败，%v", err.Error())
		return
	}

	SetValue(envProperty.Key, envProperty.Value)

	// 发布配置变更事件
	listener.PublishEvent(listener.ConfigChangeEvent{Key: envProperty.Key, Value: envProperty.Value})
}

func UpdateConfigJson(c *gin.Context) {
	valueMap := map[string]any{}
	err := isc.DataToObject(c.Request.Body, &valueMap)
	if err != nil {
		log.Printf("解析失败，%v", err.Error())
		return
	}

	key, _ := valueMap["key"]
	value, _ := valueMap["value"]

	SetValue(key.(string), value)

	// 发布配置变更事件
	listener.PublishEvent(listener.ConfigChangeEvent{Key: key.(string), Value: isc.ToJsonString(value)})
}

// 多种格式优先级：json > properties > yaml > yml
func doLoadConfigFromAbsPath(resourceAbsPath string) {
	if !strings.HasSuffix(resourceAbsPath, "/") {
		resourceAbsPath += "/"
	}
	files, err := os.ReadDir(resourceAbsPath)
	if err != nil {
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	LoadYamlFile(resourceAbsPath + "application.yaml")
	LoadYamlFile(resourceAbsPath + "application.yml")
	LoadPropertyFile(resourceAbsPath + "application.properties")
	LoadJsonFile(resourceAbsPath + "application.json")

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()
		if !strings.HasPrefix(fileName, "application") {
			continue
		}

		// 默认配置
		if fileName == "application.yaml" {
			configExist = true
			break
		} else if fileName == "application.yml" {
			configExist = true
			break
		} else if fileName == "application.properties" {
			configExist = true
			break
		} else if fileName == "application.json" {
			configExist = true
			break
		}

		profile := getActiveProfile()
		if profile != "" {
			CurrentProfile = profile
			SetValue("base.profiles.active", profile)
			currentProfile := getProfileFromFileName(fileName)
			if currentProfile == profile {
				AppendFile(resourceAbsPath + fileName)
			}
		}
	}
}

func LoadFile(filePath string) {
	extend := getFileExtension(filePath)
	extend = strings.ToLower(extend)
	if extend == "yaml" {
		configExist = true
		LoadYamlFile(filePath)
	} else if extend == "yml" {
		configExist = true
		LoadYamlFile(filePath)
	} else if extend == "properties" {
		configExist = true
		LoadPropertyFile(filePath)
	} else if extend == "json" {
		configExist = true
		LoadJsonFile(filePath)
	}
}

func AppendFile(filePath string) {
	extend := getFileExtension(filePath)
	extend = strings.ToLower(extend)
	if extend == "yaml" {
		AppendYamlFile(filePath)
	} else if extend == "yml" {
		AppendYamlFile(filePath)
	} else if extend == "properties" {
		AppendPropertyFile(filePath)
	} else if extend == "json" {
		AppendJsonFile(filePath)
	}
}

// 临时写死
// 优先级：环境变量 > 本地配置
func getActiveProfile() string {
	profile := os.Getenv("base.profiles.active")
	if profile != "" {
		return profile
	}

	profile = GetValueString("base.profiles.active")
	if profile != "" {
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
		lastIndex := strings.LastIndex(fileName, ".")
		if lastIndex > -1 {
			return fileName[lastIndex+1:]
		}
	}
	return ""
}

func LoadYamlFile(filePath string) {
	if !file.FileExists(filePath) {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		// log.Printf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	property, err := isc.YamlToProperties(string(content))
	if err != nil {
		return
	}
	valueMap, _ := isc.PropertiesToMap(property)
	appProperty.ValueMap = valueMap

	yamlMap, err := isc.YamlToMap(string(content))
	if err != nil {
		return
	}
	appProperty.ValueDeepMap = yamlMap
}

func AppendYamlFile(filePath string) {
	if !file.FileExists(filePath) {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		// log.Printf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	property, err := isc.YamlToProperties(string(content))
	if err != nil {
		return
	}
	AppendValue(property)
}

func LoadPropertyFile(filePath string) {
	if !file.FileExists(filePath) {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		// log.Printf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	valueMap, _ := isc.PropertiesToMap(string(content))
	appProperty.ValueMap = valueMap

	yamlStr, _ := isc.PropertiesToYaml(string(content))
	yamlMap, _ := isc.YamlToMap(yamlStr)
	appProperty.ValueDeepMap = yamlMap
}

func AppendPropertyFile(filePath string) {
	if !file.FileExists(filePath) {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	valueMap, err := isc.PropertiesToMap(string(content))
	if err != nil {
		return
	}
	propertiesValue, err := isc.MapToProperties(valueMap)
	if err != nil {
		return
	}

	AppendValue(propertiesValue)
}

func LoadJsonFile(filePath string) {
	if !file.FileExists(filePath) {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		// log.Printf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	yamlStr, _ := isc.JsonToYaml(string(content))
	property, _ := isc.YamlToProperties(yamlStr)
	valueMap, _ := isc.PropertiesToMap(property)
	appProperty.ValueMap = valueMap

	yamlMap, _ := isc.YamlToMap(yamlStr)
	appProperty.ValueDeepMap = yamlMap
}

func AppendJsonFile(filePath string) {
	if !file.FileExists(filePath) {
		return
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败(%v)", err)
		return
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	yamlStr, err := isc.JsonToYaml(string(content))
	if err != nil {
		return
	}
	property, err := isc.YamlToProperties(yamlStr)
	if err != nil {
		return
	}

	AppendValue(property)
}

func AppendValue(propertiesNewValue string) {
	pMap, err := isc.PropertiesToMap(propertiesNewValue)
	for k, v := range pMap {
		appProperty.ValueMap[k] = v
	}

	propertiesValueOfOriginal, err := isc.MapToProperties(appProperty.ValueMap)
	if err != nil {
		return
	}

	resultYaml, err := isc.PropertiesToYaml(propertiesValueOfOriginal)
	if err != nil {
		return
	}
	resultDeepMap, err := isc.YamlToMap(resultYaml)
	if err != nil {
		return
	}
	appProperty.ValueDeepMap = resultDeepMap
}

func SetValue(key string, value any) {
	if nil == value {
		return
	}
	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = make(map[string]interface{})
		appProperty.ValueDeepMap = make(map[string]interface{})
	} else if appProperty.ValueMap == nil {
		appProperty.ValueMap = make(map[string]interface{})
	} else if appProperty.ValueDeepMap == nil {
		appProperty.ValueDeepMap = make(map[string]interface{})
	}

	if oldValue, exist := appProperty.ValueMap[key]; exist {
		if !isc.IsBaseType(reflect.TypeOf(oldValue)) {
			if reflect.TypeOf(oldValue) != reflect.TypeOf(value) {
				return
			}
		}
	}
	propertiesValueOfOriginal, err := isc.MapToProperties(appProperty.ValueDeepMap)
	if err != nil {
		return
	}
	resultMap, err := isc.PropertiesToMap(propertiesValueOfOriginal)
	if err != nil {
		return
	}

	if reflect.ValueOf(value).Kind() == reflect.Map || reflect.ValueOf(value).Kind() == reflect.Struct {
		valueMap, err := isc.JsonToMap(isc.ObjectToJson(value))
		if err != nil {
			return
		}
		for k, v := range valueMap {
			resultMap[key + "." + k] = v
		}
	} else if reflect.ValueOf(value).Kind() == reflect.Slice || reflect.ValueOf(value).Kind() == reflect.Array {
		values := []any{}
		err := isc.DataToObject(isc.ObjectToJson(value), &values)
		if err != nil {
			return
		}
		for i, v := range values {
			resultMap[key + "[" + isc.ToString(i) + "]"] = v
		}
	}

	appProperty.ValueMap = resultMap

	mapProperties, err := isc.MapToProperties(resultMap)
	if err != nil {
		return
	}
	mapYaml, err := isc.PropertiesToYaml(mapProperties)
	if err != nil {
		return
	}
	resultDeepMap, err := isc.YamlToMap(mapYaml)
	if err != nil {
		return
	}
	appProperty.ValueDeepMap = resultDeepMap
}

func parseProperties(key string, value any, resultMap map[string]any) (map[string]any, error) {
	if reflect.ValueOf(value).Kind() == reflect.Map || reflect.ValueOf(value).Kind() == reflect.Struct {
		valueMap, err := isc.JsonToMap(isc.ObjectToJson(value))
		if err != nil {
			return resultMap, err
		}
		for k, v := range valueMap {
			resultMap[key + "." + k] = v
		}
	} else if reflect.ValueOf(value).Kind() == reflect.Slice || reflect.ValueOf(value).Kind() == reflect.Array {
		values := []any{}
		err := isc.DataToObject(isc.ObjectToJson(value), &values)
		if err != nil {
			return resultMap, err
		}
		for i, v := range values {
			resultMap[key + "[" + isc.ToString(i) + "]"] = v
			// todo 这里要添加一些东西
		}
	}
	return resultMap, nil
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

func GetValueObject(key string, targetPtrObj any) error {
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

func GetValueArray(key string) []any {
	if nil == appProperty {
		return nil
	}

	var arrayResult []any
	data := doGetValue(appProperty.ValueDeepMap, key)
	err := isc.DataToObject(data, &arrayResult)
	if err != nil {
		return arrayResult
	}
	return arrayResult
}

func GetValueArrayInt(key string) []int {
	if nil == appProperty {
		return nil
	}

	var arrayResult []int
	data := doGetValue(appProperty.ValueDeepMap, key)
	err := isc.DataToObject(data, &arrayResult)
	if err != nil {
		return arrayResult
	}
	return arrayResult
}

func GetValue(key string) any {
	if nil == appProperty {
		return nil
	}
	return doGetValue(appProperty.ValueDeepMap, key)
}

func doGetValue(parentValue any, key string) any {
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
	ValueMap     map[string]any
	ValueDeepMap map[string]any
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
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("读取文件异常(%v)", err)
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
			log.Printf("读取 application.yml 异常(%v)", err)
			return err
		}
		v1 := reflect.ValueOf(AConfig).Elem()
		o1 := v1.FieldByName("Spring").Interface()
		v2 := reflect.ValueOf(o1)
		o2 := v2.FieldByName("Profiles").Interface()
		v3 := reflect.ValueOf(o2)
		act := v3.FieldByName("Active").String()
		if act != "" && act != "default" {
			yamlAdditional, err := os.ReadFile(fmt.Sprintf("./application-%s.yml", act))
			if err != nil {
				log.Printf("读取 application-%s.yml 失败", act)
				return err
			} else {
				err = yaml.Unmarshal(yamlAdditional, AConfig)
				if err != nil {
					log.Printf("读取 application-%s.yml 异常", act)
					return err
				}
			}
		}
		return nil
	})
}
