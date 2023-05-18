## validate

validate包核查模块，用于对入参的校验

## 快速使用
这里举个例子，快速使用

### 基于isc-gobase的web的项目的示例：
```go
// main.go 文件
package main

import (
  "bytes"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/isyscore/isc-gobase/http"
    "github.com/isyscore/isc-gobase/isc"
    "github.com/isyscore/isc-gobase/logger"
    "github.com/isyscore/isc-gobase/server"
    "github.com/isyscore/isc-gobase/server/rsp"
    "github.com/isyscore/isc-gobase/validate"
    "io/ioutil"
    "strings"
)

func main() {
    server.Post("test/insert", InsertData)
    server.Run()
}

// InsertData 数据插入
func InsertData(c *gin.Context) {
    insertReq := InsertReq{}

    // 读取body数据，可以采用isc提供的工具
    err := isc.DataToObject(c.Request.Body, &insertReq)
    if err != nil {
        // ... 省略异常日志
        return
    }

    // api示例：核查入参
    if result, msg := validate.Check(insertReq); !result {
        // 参数异常
        rsp.FailedOfStandard(c, 53, msg)
        logger.Error(msg)
        return
    }

    rsp.SuccessOfStandard(c, "ok")
}

type InsertReq struct {
    Name    string `match:"value={zhou, chen}"`
    Profile string `match:"range=[0, 10)"`
}
```
```yaml
# application.yml 文件
api-module: app/sample

base:
  api:
    # api前缀
    prefix: /api
  application:
    # 应用名称
    name: sample
  server:
    # 是否启用，默认：true
    enable: true
    # 端口号
    port: 8080
    # web框架gin的配置
    gin:
      # 有三种模式：debug/release/test
      mode: release

```

请求
```shell
curl -X POST \
  http://localhost:8080/api/app/sample/test/insert \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "zhou",
	"profile": "abcde-abcde"
}'
```
返回异常
```json
{
    "code": 53,
    "data": null,
    "message": "长度不合法"
}
```

### 非web的普通类核查示例：

```go
package main

import (
    "github.com/isyscore/isc-gobase/validate"
    "github.com/isyscore/isc-gobase/logger"
)

type DemoInsert struct {
    // 对属性修饰 
    Name string `match:"value=zhou"`
    Age  int
}

func main() {
    var value DemoInsert
    var result bool
    var errMsg string

    value = DemoInsert{Name: "chen"}
    // 核查
    result, errMsg = validate.Check(value)
    if !result {
        // 属性 Name 的值 chen 不在只可用列表 [zhou] 中 
        logger.Error(errMsg)
    }
}
```
说明：<br/>

1. 这里提供方法Check，用于核查是否符合条件
2. 提供标签match，标签内容中提供匹配器：value，该匹配器表示匹配的具体的一些值

## api说明
只提供两个Api，`Check` 和 `CheckWithParameter`
```go
// 入参：
//  @any            待核查对象
//  @fieldNames     待核查对象的核查属性；不指定则核查所有属性名
// 返回值：
//  bool    是否匹配合法：true-合法，false-不合法
//  string  不合法对应的说明
func Check(object any, fieldNames ...string) (bool, string) {}

// 入参：
//  @parameterMap   额外的参数，用于在自定义函数中进行使用，可见下面的customize的用法
//  @any            待核查对象
//  @fieldNames     待核查对象的核查属性；不指定则核查所有属性名
// 返回值：
//  bool    是否匹配合法：true-合法，false-不合法
//  string  不合法对应的说明
func CheckWithParameter(parameterMap map[string]interface{}, object interface{}, fieldNames ...string) (bool, string) {}
```

## 更多功能

这里将核查部分分为匹配和处理两部分，匹配可以有多种的匹配器，核查的逻辑是只要有任何一个匹配器匹配上则认为匹配上，处理模块用于对匹配上的结果进行处理，比如返回指定的异常，或者匹配后接受还是拒绝对应的值，或者匹配后将某个值更改掉（待支持）

#### 匹配模块

- value：匹配指定的值
- isBlank：值是否为空字符
- isUnBlank：值是否为非空字符
- range：匹配数值的范围（最大值和最小值，用法是数学表达式）：数值（整数和浮点数）的大小、字符串的长度、数组的长度、时间的范围、时间的移动
- model：匹配指定的类型：
    - id_card：身份证
    - phone: 手机号
    - fixed_phone:固定电话
    - mail: 邮件地址
    - ip: ip地址
- condition：修饰的属性的表达式的匹配，提供#current和#root占位符，用于获取相邻属性的值
- regex：匹配正则表达式
- customize：匹配自定义的回调函数

#### 处理模块

- errMsg: 自定义的异常
- accept: 匹配后接受还是拒绝
- disable: 是否启用匹配，默认启用


## 匹配模块

匹配器可以有多个一起修饰，只要匹配上一个，则认为匹配上

### 1. 值匹配器：value
匹配指定的一些值，可以修饰一个，也可以修饰多个值，可以修饰字符，也可修饰整数（int、int8、int16、int32、int64）、无符号整数（uint、uint8、uint16、uint32、uint64）、浮点数（float32、float64）、bool类型和string类型。<br/>

提示：
 - 中间逗号也可以为中文，为了防止某些手误写错为中文字符


```go
// 修饰一个值
type ValueBaseEntityOne struct {
    Name string `match:"value=zhou"`
    Age  int    `match:"value=12"`
}

// 修饰一个值
type ValueBaseEntity struct {
    Name string `match:"value={zhou, 宋江}"`
    Age  int    `match:"value={12, 13}"`
}
```

如果有自定义类型嵌套，则可以使用标签`check`，用于解析复杂结构
```go
type ValueInnerEntity struct {
    InnerName string `match:"value={inner_zhou, inner_宋江}"`
    InnerAge  int    `match:"value={2212, 2213}"`
}

type ValueStructEntity struct {
    Name string `match:"value={zhou, 宋江}"`
    Age  int    `match:"value={12, 13}"`

    Inner ValueInnerEntity `match:"check"`
}
```
修饰的结构可以有如下
- 自定义结构
- 数组/分片：对应类型只有为复杂结构才会核查
- map：其中的key和value类型只有是复杂结构才会核查

### 2. 空值匹配器：isBlank
匹配string类型的值是否为空字符，false：字符不为空则匹配上，true：字符为空则匹配上
```go
// 默认为true
type IsBlankEntity2 struct {
    Name string `match:"isBlank"`
    Age  int
}

// 同上
type IsBlankEntity3 struct {
    Name string `match:"isBlank=true"`
    Age  int
}

type IsBlankEntity1 struct {
    Name string `match:"isBlank=false"`
    Age  int
}

```

### 3. 非空匹配器：isUnBlank
匹配string类型的值是否为非空字符，true：字符非空则匹配上，false：字符为空则匹配上
```go
// 默认为true
type IsBlankEntity2 struct {
    Name string `match:"isBlank"`
    Age  int
}

// 同上
type IsBlankEntity3 struct {
    Name string `match:"isBlank=true"`
    Age  int
}

type IsBlankEntity1 struct {
    Name string `match:"isBlank=false"`
    Age  int
}

```

### 4. 范围匹配器：range
匹配类型的指定范围，方式使用数学表达式"["、"]"、"("、")"，使用数学表达式的开闭符号
- [：表示大于等于
- ]：表示小于等于
- (：表示大于
- )：表示小于

比如：<br/>
- [1, 10)：表示大于等于1而且小于10
- (1, 10)：表示大于1而且小于10
- (1,)：表示大于1
- (,10]：表示小于等于10，也可以[,10]

修饰的类型有
- 整数：比较大小，int、int8、int16、int32、int64
- 无符号整数：比较大小，uint、uint8、uint16、uint32、uint64
- 浮点数：比较大小，float32、float64
- 分片：匹配分片的长度
- 字符串：匹配字符串的长度
- 时间类型（time.Time）：时间的范围，时间格式支持如下
  - yyyy
  - yyyy-MM
  - yyyy-MM-dd
  - yyyy-MM-dd HH
  - yyyy-MM-dd HH:mm
  - yyyy-MM-dd HH:mm:ss
  - yyyy-MM-dd HH:mm:ss.SSS
  - now：表示当前时间
- 变量时间：除了使用数学表达式外，还支持past、future这两个关键字，用于表示过去和未来的时间
- 时间计算：用于计算当前的时间向前或者向后推几个小时或者几分钟这种；(-1M, )：表示最近一个月的时间；(-1M3d, )：表示最近一个月零3天的时间，表示大于时间向前推一个月零三天的时间<br/>
  - -/+：表示往前推还是往后推
  - y：年
  - M：月
  - d：日
  - h：小时
  - m：分钟
  - s：秒

```go
// 整数类型1
type RangeIntEntity1 struct {
    Age  int `match:"range=[1, 2]"`
}

// 整数类型2
type RangeIntEntity2 struct {
    Age  int `match:"range=[3，]"`
}

// 整数类型3
type RangeIntEntity3 struct {
    Age  int `match:"range=[3,)"`
}

// 浮点数类型
type RangeFloatEntity struct {
    Money float32 `match:"range=[10.37， 20.31]"`
}

// 字符类型
type RangeStringEntity struct {
    Name string `match:"range=[2, 12]"`
}

// 分片类型
type RangeSliceEntity struct {
    Age  []int `match:"range=[2, 6]"`
}

// 时间类型1
type RangeTimeEntity1 struct {
    CreateTime time.Time `match:"range=[2019-07-13 12:00:23.321, 2019-08-23 12:00:23.321]"`
}

// 时间类型2
type RangeTimeEntity2 struct {
    CreateTime time.Time `match:"range=[2019-07-13 12:00:23.321, ]"`
}

// 时间类型3
type RangeTimeEntity3 struct {
    CreateTime time.Time `match:"range=(, 2019-07-23 12:00:23.321]"`
}

// 时间类型4
type RangeTimeEntity4 struct {
    CreateTime time.Time `match:"range=[2019-07-23 12:00:23.321, now)"`
}

// 时间类型4
type RangeTimeEntity5 struct {
    CreateTime time.Time `match:"range=past"`
}

// 时间类型4
type RangeTimeEntity6 struct {
    CreateTime time.Time `match:"range=future"`
}

// 时间计算：年
type RangeTimeCalEntity1 struct {
    CreateTime time.Time `match:"range=(-1y, )"`
}

// 时间计算：月
type RangeTimeCalEntity2 struct {
    CreateTime time.Time `match:"range=(-1M, )"`
}

// 时间计算：月日
// 顺序不能乱：yMdhms，中间可以为空，比如：-1y3d
type RangeTimeCalEntity2And1 struct {
    CreateTime time.Time `match:"range=(-1M3d, )"`
}

// 时间计算：日
type RangeTimeCalEntity3 struct {
    CreateTime time.Time `match:"range=(-3d, )"`
}

// 时间计算：时
type RangeTimeCalEntity4 struct {
    CreateTime time.Time `match:"range=(-4h, )"`
}

// 时间计算：分
type RangeTimeCalEntity5 struct {
    CreateTime time.Time `match:"range=(-12m, )"`
}

// 时间计算：秒
type RangeTimeCalEntity6 struct {
    CreateTime time.Time `match:"range=(-120s, )"`
}

// 时间计算：正负号
type RangeTimeCalEntity7 struct {
    CreateTime time.Time `match:"range=(2h, )"`
}
```

### 5. 类型匹配器：model
类型匹配器：指定的几种内置类型进行匹配
- id_card：身份证
- phone: 手机号
- fixed_phone:固定电话
- mail: 邮件地址
- ip: ip地址

```go
type ValueModelIdCardEntity struct {
    Data string `match:"model=id_card"`
}

type ValueModelPhone struct {
    Data string `match:"model=phone"`
}

type ValueModelFixedPhoneEntity struct {
    Data string `match:"model=fixed_phone"`
}

type ValueModelEmailEntity struct {
    Data string `match:"model=mail"`
}

type ValueModelIpAddressEntity struct {
    Data string `match:"model=ip"`
}
```

### 6. 表达式匹配器：condition
表达式匹配器：用于数学计算表达式进行计算，表达式是返回bool类型的表达式。提供两个占位符
- \#current：当前修饰的值
- \#root：当前属性所在的对象，比如：#root.Age，表示当前对象中的其他属性Age的值

```go
// 测试基本表达式
type ValueConditionEntity1 struct {
    Data1 int `match:"condition=#current + #root.Data2 > 100"`
    Data2 int `match:"condition=#current < 20"`
    Data3 int `match:"condition=(++#current) >31"`
}

// 测试表达式
type ValueConditionEntity2 struct {
    Age   int `match:"condition=#root.Judge"`
    Judge bool
}
```

### 7. 正则表达式匹配器：regex
正则表达式匹配器：用于匹配自定义的正则表达式

```go
type ValueRegexEntity struct {
    Name string `match:"regex=^zhou.*zhen$"`
    Age  int    `match:"regex=^\\d+$"`
}
```

### 8. 自定义回调匹配器：customize
该匹配器可以用于自定义扩展，比如实际业务场景，某个字段在数据库中存在，这种情况就需要用户自定义扩展

比如：
```go
package fun

type CustomizeEntity1 struct {
    // fun.Judge1是对应的函数
    Name string `match:"customize=judge1Name"`
}

func JudgeString1(name string) bool {
    if name == "zhou" || name == "宋江" {
        return true
    }

    return false
}

// 由于go反射功能没那么强，因此需要用户自己先将函数和name进行注册
func init() {
  validate.RegisterCustomize("judge1Name", JudgeString1)
}
```
##### 说明：

其中自定义的函数有相关的要求，参数可以为一个，也可以为两个，也可以为三个<br/>
其中的参数类型有严格限制
- 属性类型
- 属性所在对象类型
- 外部参数类型`map[string]interface{}`


参数：<br/>
- 一个参数：
  - 1：属性类型
  - 2：属性所在对象类型
- 两个参数
  - 1：属性所在对象类型，2：属性类型
  - 1：属性所在对象类型，2：外部参数类型`map[string]interface{}`
  - 1：属性类型，2：属性所在对象类型
  - 1：属性类型，2：外部参数类型`map[string]interface{}`
- 三个参数：前两个参数为：属性类型和所在对象类型的组合
  - 1：属性类型，2：属性所在对象类型，3：外部参数类型`map[string]interface{}`
  - 1：属性所在对象类型，2：属性类型，3：外部参数类型`map[string]interface{}`

返回值：<br/>
- 一个值：则为bool类型（表示是否匹配上）
- 两个值：第一个为bool类型（表示是否匹配上），第二个为string类型（匹配或者没有匹配上的自定义错误）

更多详情请见测试类`customize_test.go` <br/><br/>
示例：
```go
package fun

import (
    "fmt"
    "github.com/isyscore/isc-gobase/validate"
)

type CustomizeEntity2 struct {
    Name string `match:"customize=judge2Name"`
}

type CustomizeEntity3 struct {
    Name string `match:"customize=judge3Name"`
    Age  int
}

func JudgeString2(name string) (bool, string) {
    if name == "zhou" || name == "宋江" {
        return true, ""
    }

    return false, "没有命中可用的值'zhou'和'宋江'"
}

func JudgeString3(customize CustomizeEntity3, name string) (bool, string) {
    if name == "zhou" || name == "宋江" {
    if customize.Age > 12 {
        return true, ""
    } else {
        return false, "用户[" + name + "]" + "没有满足年龄age > 12，" + "当前年龄为：" + fmt.Sprintf("%v", customize.Age)
    }

    } else {
        return false, "没有命中可用的值'zhou'和'宋江'"
    }
}

// 由于go反射功能没那么强，因此需要用户自己先将函数和name进行注册
func init() {
    validate.RegisterCustomize("judge2Name", JudgeString2)
    validate.RegisterCustomize("judge3Name", JudgeString3)
}
```

## 处理模块
当匹配后如何处理，这里分为了如下几种处理

### 1. 匹配上接受/拒绝：accept
匹配后是接收还是拒绝，目前业内的处理方式都是匹配后接收，在概念上其实叫白名单，对于黑名单的处理，业内是没有，而我们这里用accept实现白名单和黑名单的概念

```go
type AcceptEntity1 struct {
    // 表示只要匹配上name为zhou的，则拒绝
    Name string `match:"value=zhou" accept:"false"`
    Age  int
}

type AcceptEntity2 struct {
    // 表示name为空，则拒绝，表示需要非空才行，以下同`match:"isUnBlank"`
    Name string `match:"isBlank" accept:"false"`
    Age  int
}

// 只要任何一个匹配上，则接受
type AcceptEntity3 struct {
    Name string `match:"isBlank=true value=zhou" accept:"true"`
    Age  int
}
```

### 2. 自定义异常：errMsg
匹配后返回自定义的异常，提供了两个占位符#current表示修饰的当前属性的值，#root当前属性所在的结构的值，#root.Age表示当前结构中的属性Age对应的值

```go
type ErrMsgEntity1 struct {
    Name string `match:"value=zhou" errMsg:"对应的值不合法"`
    Age  int
}

type ErrMsgEntity2 struct {
    Name string `match:"value=zhou"`
    Age  int    `match:"condition=#current > 10" errMsg:"当前的值不合法，应该大于10，当前值为#current，对应的名字为#root.Name"`
}
```

### 3. 启用：disable
表示是否启用整个核查
```go
type DisableEntity1 struct {
    Name string `match:"value=zhou" disable:"true"`
    Age  int
}

```
