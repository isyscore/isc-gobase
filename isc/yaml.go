package isc

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

/**
 *  1.yaml <---> properties
 *  2.yaml <---> json
 *  3.yaml <---> map
 *  4.yaml <---> list
 *  5.yaml <---> kvList
 */

// NewLine 换行符
var NewLine = "\n"

// IndentBlanks 缩进空格
var IndentBlanks = "  "

// SignSemicolon 分号连接符
var SignSemicolon = ":"

// SignEqual 等号连接符
var SignEqual = "="

// Dot 点
var Dot = "."

// ArrayBlanks 数组缩进
var ArrayBlanks = "- "

// NewLineDom yaml的value换行符
var YamlNewLineDom = "|\n"

var rangePattern = regexp.MustCompile("^(.*)\\[(\\d*)\\]$")

type TypeEnum int8

const (
	YAML       TypeEnum = 0
	PROPERTIES TypeEnum = 1
	JSON       TypeEnum = 2
	STRING     TypeEnum = 3
)

type Properties struct {
	Value map[string]string
}

type StringPair struct {
	Left  string
	Right string
}

type ConvertError struct {
	errMsg string
}

func (convertError *ConvertError) Error() string {
	return convertError.errMsg
}

func IsYaml(content string) bool {
	if !strings.Contains(content, ":") && !strings.Contains(content, "-") {
		return false
	}

	_, err := YamlToProperties(content)
	if err != nil {
		return false
	}
	return true
}

func IsProperty(content string) bool {
	if !strings.Contains(content, "=") {
		return false
	}

	yml, _ := PropertiesToYaml(content)
	_, err := YamlToProperties(yml)
	if err != nil {
		return false
	}
	return true
}

func IsJson(content string) bool {
	if !strings.HasPrefix(content, "{") && !strings.HasPrefix(content, "[") {
		return false
	}

	var object interface{}
	err := json.Unmarshal([]byte(content), &object)
	if err != nil {
		return false
	}
	return true
}

func YamlToProperties(contentOfYaml string) (string, error) {
	// yaml 到 map
	dataMap, err := YamlToMap(contentOfYaml)
	if err != nil {
		log.Printf("YamlToPropertiesStr error: %v, content: %v", err, contentOfYaml)
		return "", err
	}

	return MapToProperties(dataMap)
}

func YamlToPropertiesWithKey(key string, contentOfYaml string) (string, error) {
	if "" == key {
		return "", &ConvertError{errMsg: "key is nil"}
	}
	if "" == contentOfYaml {
		return key + "=\"\"", nil
	}
	if !strings.Contains(contentOfYaml, ":") && !strings.Contains(contentOfYaml, "-") {
		return "", &ConvertError{errMsg: "content is illegal for yaml"}
	}

	contentOfYaml = strings.TrimSpace(contentOfYaml)
	if strings.HasPrefix(contentOfYaml, "-") {
		dataMap := map[string]interface{}{}
		kvList, err := YamlToList(contentOfYaml)
		if err != nil {
			log.Printf("YamlToPropertiesWithKey error: %v, content: %v", err, contentOfYaml)
			return "", err
		}

		dataMap[key] = kvList
		yamlStr, err := ObjectToYaml(dataMap)
		if err != nil {
			return "", err
		}
		return YamlToProperties(yamlStr)
	}

	property, err := YamlToProperties(contentOfYaml)
	if err != nil {
		log.Printf("YamlToPropertiesWithKey error: %v, content: %v", err, contentOfYaml)
		return "", err
	}

	return propertiesAppendPrefixKey(key, property)
}

func JsonToMap(contentOfJson string) (map[string]interface{}, error) {
	if contentOfJson == "" || strings.HasPrefix(contentOfJson, "{") || strings.HasPrefix(contentOfJson, "[") {
		return nil, &ChangeError{ErrMsg: "不符合json格式"}
	}
	resultMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(contentOfJson), &resultMap)
	if err != nil {
		log.Printf("JsonToMap, error: %v, content: %v", err, contentOfJson)
		return nil, err
	}

	return resultMap, nil
}

func YamlToMap(contentOfYaml string) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(contentOfYaml), &resultMap)
	if err != nil {
		log.Printf("YamlToMap, error: %v, content: %v", err, contentOfYaml)
		return nil, err
	}

	return resultMap, nil
}

func YamlToJson(contentOfYaml string) (string, error) {
	if contentOfYaml != "-" && strings.Contains(contentOfYaml, ":") {
		return "", &ConvertError{errMsg: "the content is invalidate for json"}
	}

	var data interface{}
	err := yaml.Unmarshal([]byte(contentOfYaml), &data)
	if err != nil {
		log.Printf("YamlToList, error: %v, content: %v", err, contentOfYaml)
		return "", err
	}

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func YamlToKvList(contentOfYaml string) ([]StringPair, error) {
	if !strings.Contains(contentOfYaml, ":") && !strings.Contains(contentOfYaml, "-") {
		return nil, nil
	}

	property, err := YamlToProperties(contentOfYaml)
	if err != nil {
		return nil, err
	}

	propertiesLineWordList := GetPropertiesItemLineList(property)
	var pairs []StringPair
	for _, element := range propertiesLineWordList {
		element = strings.TrimSpace(element)
		if "" == element {
			continue
		}
		values := strings.SplitN(element, SignEqual, 2)
		pairs = append(pairs, StringPair{Left: values[0], Right: values[1]})
	}

	return pairs, nil
}

func YamlToList(contentOfYaml string) ([]interface{}, error) {
	if !strings.HasPrefix(strings.TrimSpace(contentOfYaml), "-") {
		return []interface{}{}, &ConvertError{errMsg: "the content of yaml not start with '-'"}
	}
	var resultList []interface{}
	err := yaml.Unmarshal([]byte(contentOfYaml), &resultList)
	if err != nil {
		log.Printf("YamlToList, error: %v, content: %v", err, contentOfYaml)
		return nil, err
	}

	return resultList, nil
}

func YamlCheck(content string) error {
	if "" == content {
		return &ConvertError{errMsg: "yaml is empty"}
	}

	if !strings.Contains(content, ":") && !strings.Contains(content, "-") {
		return &ConvertError{errMsg: "yaml content does not contain ':' nor '-'"}
	}

	if strings.Contains(content, "---\n") {
		return &ConvertError{errMsg: "yaml not support import for content of '---'"}
	}

	_, err := YamlToProperties(content)
	if err != nil {
		return err
	}
	return nil
}

func PropertiesToMap(contentOfProperties string) (map[string]interface{}, error) {
	if !strings.Contains(contentOfProperties, "=") {
		return nil, &ConvertError{errMsg: "the content is illegal for properties"}
	}

	var resultMap = make(map[string]interface{})
	propertiesLineWordList := GetPropertiesItemLineList(contentOfProperties)
	for _, line := range propertiesLineWordList {
		line = strings.TrimSpace(line)
		if "" == line {
			continue
		}

		lineKVs := strings.SplitN(line, "=", 2)
		key := lineKVs[0]
		value := lineKVs[1]

		if strings.Contains(value, "\n") {
			value = YamlNewLineDom + value
		}

		resultMap[key] = value
	}

	return resultMap, nil
}

func propertiesAppendPrefixKey(key string, propertiesContent string) (string, error) {
	itemLines := GetPropertiesItemLineList(propertiesContent)
	var datas []string
	for _, line := range itemLines {
		if !strings.Contains(line, SignEqual) {
			continue
		}

		kvs := strings.SplitN(line, SignEqual, 2)
		datas = append(datas, key+Dot+kvs[0]+SignEqual+kvs[1])
	}

	return strings.Join(datas, NewLine), nil
}

func deepPut(dataMap map[string]interface{}, key string, value interface{}) map[string]interface{} {
	mapValue, exist := dataMap[key]
	if !exist {
		dataMap[key] = value
	} else {
		if reflect.Map == reflect.TypeOf(value).Kind() {
			leftMap := mapValue.(map[string]interface{})
			rightMap := value.(map[string]interface{})

			for rightMapKey := range rightMap {
				leftMap = deepPut(leftMap, rightMapKey, rightMap[rightMapKey])
			}
			dataMap[key] = leftMap
		}
	}

	return dataMap
}

func PropertiesToYaml(contentOfProperties string) (string, error) {
	var yamlLineList []string
	var yamlNodes []YamlNode
	propertiesLineWordList := GetPropertiesItemLineList(contentOfProperties)
	for _, line := range propertiesLineWordList {
		line = strings.TrimSpace(line)
		if line != "" {
			// 注释数据不要
			if strings.HasPrefix(line, "#") {
				continue
			}

			index := strings.Index(line, SignEqual)
			if index > -1 {
				key := line[:index]
				value := line[index+1:]

				if strings.Contains(value, "\n") {
					value = YamlNewLineDom + value
				}

				lineWordList := strings.Split(key, ".")
				lineWordList, yamlNodes = wordToNode(lineWordList, yamlNodes, nil, false, -1, appendSpaceForArrayValue(value))
			}
		}
	}
	yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNodes, false, "")
	return strings.Join(yamlLineList, "\n") + "\n", nil
}

func ObjectToYaml(value interface{}) (string, error) {
	bytes2, err := yaml.Marshal(value)
	if err != nil {
		log.Printf("ObjectToYaml error: %v, content: %v", err, value)
		return "", &ConvertError{errMsg: "ObjectToYaml error"}
	}
	return string(bytes2), nil
}

// 进行深层嵌套的map数据处理
func MapToProperties(dataMap map[string]interface{}) (string, error) {
	var propertyStrList []string
	for key, value := range dataMap {
		valueKind := reflect.TypeOf(value).Kind()
		switch valueKind {
		case reflect.Map:
			{
				propertyStrList = doMapToProperties(propertyStrList, value, prefixWithDOT("")+key)
			}
		case reflect.Array, reflect.Slice:
			{
				objectValue := reflect.ValueOf(value)
				for index := 0; index < objectValue.Len(); index++ {
					propertyStrList = doMapToProperties(propertyStrList, objectValue.Index(index).Interface(), prefixWithDOT("")+key+"["+strconv.Itoa(index)+"]")
				}
			}
		case reflect.String:
			objectValue := reflect.ValueOf(value)
			objectValueStr := strings.ReplaceAll(objectValue.String(), "\n", "\\\n")
			propertyStrList = append(propertyStrList, prefixWithDOT("")+key+SignEqual+objectValueStr)
		default:
			propertyStrList = append(propertyStrList, prefixWithDOT("")+key+SignEqual+fmt.Sprintf("%v", value))
		}
	}
	resultStr := ""
	for _, propertyStr := range propertyStrList {
		resultStr += propertyStr + "\n"
	}

	return resultStr, nil
}

func KvToProperties(key, value string, valueType TypeEnum) (string, error) {
	switch valueType {
	case YAML:
		return YamlToPropertiesWithKey(key, value)
	case JSON:
		value, err := JsonToYaml(value)
		if err != nil {
			return "", err
		}
		return YamlToPropertiesWithKey(key, value)
	case PROPERTIES:
		return propertiesAppendPrefixKey(key, value)
	case STRING:
		return key + "=" + appendSpaceForArrayValue(value), nil
	default:
		break
	}
	return "", nil
}

func JsonToYaml(contentOfJson string) (string, error) {
	if !strings.HasPrefix(contentOfJson, "{") && !strings.HasPrefix(contentOfJson, "[") {
		return "", &ConvertError{errMsg: "content is not json"}
	}

	if "[]" == contentOfJson {
		return "", nil
	}

	// go中的json转换不识别"'"字符，因此这里将其转换为"""这种字符
	contentOfJsonTem := strings.ReplaceAll(contentOfJson, "'", "\"")

	var object interface{}
	err := json.Unmarshal([]byte(contentOfJsonTem), &object)
	if err != nil {
		return "", err
	}

	return ObjectToYaml(object)
}

func PropertiesEntityToYaml(properties Properties) (string, error) {
	if properties.Value == nil {
		return "", &ConvertError{"PropertiesEntityToYaml value is empty"}
	}

	var content = ""
	for key, value := range properties.Value {
		content += key + "=" + value + "\n"
	}
	return PropertiesToYaml(content)
}

func GetPropertiesItemLineList(content string) []string {
	if "" == content {
		return []string{}
	}

	lineList := strings.Split(content, NewLine)
	var itemLineList []string
	var stringAppender string
	for _, line := range lineList {
		if strings.HasSuffix(content, "\\") {
			stringAppender += line + "\n"
		} else {
			stringAppender += line
			itemLineList = append(itemLineList, stringAppender)
			stringAppender = ""
		}
	}
	return itemLineList
}

func formatPropertiesToYaml(yamlLineList []string, yamlNodes []YamlNode, lastNodeArrayFlag bool, blanks string) []string {
	var beforeNodeIndex = -1
	var equalSign string

	for _, yamlNode := range yamlNodes {
		value := yamlNode.value

		equalSign = SignSemicolon
		if "" != value {
			equalSign = SignSemicolon + " "
		}

		yamlNode.resortValue()

		name := yamlNode.name
		if lastNodeArrayFlag {
			if "" == name {
				yamlLineList = append(yamlLineList, blanks+ArrayBlanks+stringValueWrap(value))
			} else {
				if -1 != beforeNodeIndex && beforeNodeIndex == yamlNode.lastNodeIndex {
					yamlLineList = append(yamlLineList, blanks+IndentBlanks+name+equalSign+stringValueWrap(value))
				} else {
					yamlLineList = append(yamlLineList, blanks+ArrayBlanks+name+equalSign+stringValueWrap(value))
				}
			}
			beforeNodeIndex = yamlNode.lastNodeIndex
		} else {
			yamlLineList = append(yamlLineList, blanks+name+equalSign+stringValueWrap(value))
		}

		if yamlNode.arrayFlag {
			if lastNodeArrayFlag {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.valueList, true, IndentBlanks+IndentBlanks+blanks)
			} else {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.valueList, true, IndentBlanks+blanks)
			}
		} else {
			if lastNodeArrayFlag {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.children, false, IndentBlanks+IndentBlanks+blanks)
			} else {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.children, false, IndentBlanks+blanks)
			}
		}
	}
	return yamlLineList
}

func wordToNode(lineWordList []string, nodeList []YamlNode, parentNode *YamlNode, lastNodeArrayFlag bool, index int, value string) ([]string, []YamlNode) {
	if len(lineWordList) == 0 {
		if lastNodeArrayFlag {
			node := YamlNode{value: value, lastNodeIndex: -1}
			nodeList = append(nodeList, node)
		}
	} else {
		nodeName := lineWordList[0]
		nodeName, nextIndex := peelArray(nodeName)

		var node YamlNode
		if nil != parentNode {
			node = YamlNode{name: nodeName, parent: parentNode, lastNodeIndex: index}
		} else {
			node = YamlNode{name: nodeName, lastNodeIndex: index}
		}
		lineWordList = lineWordList[1:]

		//如果节点下面的子节点数量为0，则为终端节点，也就是赋值节点
		if len(lineWordList) == 0 {
			if -1 == nextIndex {
				node.value = value
			}
		}

		// nextIndex 不空，表示当前节点为数组，则之后的数据为他的节点数据
		if nextIndex != -1 {
			node.arrayFlag = true
			var hasEqualsName = false

			//遍历查询节点是否存在
			for innerIndex := range nodeList {
				//如果节点名称已存在，则递归添加剩下的数据节点
				if nodeName == nodeList[innerIndex].name && nodeList[innerIndex].arrayFlag {
					yamlNodeIndex := nodeList[innerIndex].lastNodeIndex
					if -1 == yamlNodeIndex || index == yamlNodeIndex {
						hasEqualsName = true
						lineWordList, nodeList[innerIndex].valueList = wordToNode(lineWordList, nodeList[innerIndex].valueList, node.parent, true, nextIndex, appendSpaceForArrayValue(value))
					}
				}
			}

			//如果遍历结果为节点名称不存在，则递归添加剩下的数据节点，并把新节点添加到上级yamlTree的子节点中
			if !hasEqualsName {
				lineWordList, node.valueList = wordToNode(lineWordList, node.valueList, node.parent, true, nextIndex, appendSpaceForArrayValue(value))
				nodeList = append(nodeList, node)
			}
		} else {
			var hasEqualsName = false
			for innerIndex := range nodeList {
				if !lastNodeArrayFlag {
					//如果节点名称已存在，则递归添加剩下的数据节点
					if nodeName == nodeList[innerIndex].name {
						hasEqualsName = true
						lineWordList, nodeList[innerIndex].children = wordToNode(lineWordList, nodeList[innerIndex].children, &nodeList[innerIndex], false, nextIndex, appendSpaceForArrayValue(value))
					}
				} else {
					//如果节点名称已存在，则递归添加剩下的数据节点
					if nodeName == nodeList[innerIndex].name {
						yamlNodeIndex := nodeList[innerIndex].lastNodeIndex
						if -1 == yamlNodeIndex || index == yamlNodeIndex {
							hasEqualsName = true
							lineWordList, nodeList[innerIndex].children = wordToNode(lineWordList, nodeList[innerIndex].children, &nodeList[innerIndex], true, nextIndex, appendSpaceForArrayValue(value))
						}
					}
				}
			}

			//如果遍历结果为节点名称不存在，则递归添加剩下的数据节点，并把新节点添加到上级yamlTree的子节点中
			if !hasEqualsName {
				lineWordList, node.children = wordToNode(lineWordList, node.children, &node, false, nextIndex, appendSpaceForArrayValue(value))
				nodeList = append(nodeList, node)
			}
		}
	}
	return lineWordList, nodeList
}

func doMapToProperties(propertyStrList []string, value interface{}, prefix string) []string {
	valueKind := reflect.TypeOf(value).Kind()
	switch valueKind {
	case reflect.Map:
		{
			// map结构
			if reflect.ValueOf(value).Len() == 0 {
				return []string{}
			}

			for mapR := reflect.ValueOf(value).MapRange(); mapR.Next(); {
				mapKey := mapR.Key().Interface()
				mapValue := mapR.Value().Interface()
				propertyStrList = doMapToProperties(propertyStrList, mapValue, prefixWithDOT(prefix)+fmt.Sprintf("%v", mapKey))
			}
		}
	case reflect.Array, reflect.Slice:
		{
			objectValue := reflect.ValueOf(value)
			for index := 0; index < objectValue.Len(); index++ {
				propertyStrList = doMapToProperties(propertyStrList, objectValue.Index(index).Interface(), prefix+"["+strconv.Itoa(index)+"]")
			}
		}
	case reflect.String:
		objectValue := reflect.ValueOf(value)
		objectValueStr := strings.ReplaceAll(objectValue.String(), "\n", "\\\n")
		propertyStrList = append(propertyStrList, prefix+SignEqual+objectValueStr)
	default:
		objectValue := fmt.Sprintf("%v", reflect.ValueOf(value))
		propertyStrList = append(propertyStrList, prefix+SignEqual+objectValue)
	}
	return propertyStrList
}

func prefixWithDOT(prefix string) string {
	if "" == prefix {
		return ""
	}
	return prefix + "."
}

func peelArray(nodeName string) (string, int) {
	var index = -1
	var name = nodeName
	var err error

	subData := rangePattern.FindAllStringSubmatch(nodeName, -1)
	if len(subData) > 0 {
		name = subData[0][1]
		indexStr := subData[0][2]
		if "" != indexStr {
			index, err = strconv.Atoi(indexStr)
			if err != nil {
				log.Printf("解析错误, nodeName=" + nodeName)
				return "", -1
			}
		}
	}
	return name, index
}

//
// 将yaml对应的这种value进行添加前缀空格，其中value为key1对应的value
// test:
//   key1: |
//     value1
//     value2
//     value3
// 对应的值
// {@code
// |
// value1
// value2
// value3
// }
//
// @param Value 待转换的值比如{@code
//              test:
//              key1: |
//              value1
//              value2
//              value3
//              }
// @return 添加前缀空格之后的处理
// {@code
// |
// value1
// value2
// value3
// }
//
func appendSpaceForArrayValue(value string) string {
	if !strings.HasPrefix(value, YamlNewLineDom) {
		return value
	}

	value = value[len(YamlNewLineDom):]
	valueTems := strings.Split(value, "\\n")

	strs := []string{}
	for _, element := range valueTems {
		tem := element
		if strings.HasSuffix(element, "\\") {
			tem = element[:len(element)-1]
		}
		strs = append(strs, IndentBlanks+tem)
	}
	return YamlNewLineDom + strings.Join(strs, "\n")
}

func stringValueWrap(value string) string {
	if "" == value {
		return ""
	}
	// 对数组的数据进行特殊处理
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return "'" + value + "'"
	}
	return value
}

type YamlNode struct {

	// 父节点
	parent *YamlNode

	// 只有parent为null时候，该值才可能有值
	projectRemark string

	// name
	name string

	// value
	value string

	// 子节点
	children []YamlNode

	// 数组标示
	arrayFlag bool

	// 存储的数组中的前一个节点的下标
	lastNodeIndex int

	// 只有数组标示为true，下面的value才有值
	valueList []YamlNode
}

func (yamlNode *YamlNode) resortValue() {
	if !yamlNode.arrayFlag || len(yamlNode.valueList) == 0 {
		return
	}

	sort.Slice(yamlNode.valueList, func(i, j int) bool {
		a := yamlNode.valueList[i]
		b := yamlNode.valueList[j]

		if -1 == a.lastNodeIndex || -1 == b.lastNodeIndex {
			return false
		}
		return a.lastNodeIndex < b.lastNodeIndex
	})

	for _, node := range yamlNode.valueList {
		node.resortValue()
	}
}
