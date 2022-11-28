## bean
对象进行管理

功能：
- 统一管理对象
- 动态的查看对应对象的属性
- 动态的调用对应对象的方法

场景：适合非业务开发的工具中，比如线上调试工具

### api
```go
// 注册对象；注意：对象只可为指针类型
func AddBean(beanName string, beanPtr any) {}

// 对象存在否
func ExistBean(beanName string) bool {}

// 获取对象
func GetBean(beanName string) any {}

// 获取对象key
func GetBeanNames(beanName string) []string {}

// 查看：对象属性值
func GetField(beanName, fieldName string) any {}

// 修改：对象属性值
func SetField(beanName, fieldName string) any {}

// 执行对象的函数
func CallFun(beanName, methodName string, parameterValueMap map[string]any) []any {}
```

### 示例
```go
func TestAddBean(t *testing.T) {
    tt := TestEntity{Name: "hello", Age: 12}
    // 注册
    bean.AddBean("test", &tt)

    // 获取bean
    t1 := bean.GetBean("test")
    t2 := t1.(*TestEntity)
    assert.Equal(t, t2.Name, tt.Name)
}

func TestGetFieldShow(t *testing.T) {
    tt := TestEntity{Name: "value"}
    // 添加注册
    bean.AddBean("test", &tt)

    // 获取值
    actValue := bean.GetField("test", "Name")
    assert.Equal(t, actValue, "value")

    // 修改值
    bean.SetField("test", "Name", "value-change")

    // 查看
    actValue = bean.GetField("test", "Name")
    assert.Equal(t, actValue, "value-change")
}

func TestCallFun1(t *testing.T) {
    tt := TestEntity{}

    // 添加bean
    bean.AddBean("test", &tt)

    parameterMap := map[string]any{}
    // 说明：参数map中的key只可为p1、p2、p3...，用于表示参数的顺序
    aparameterMap["p1"] = "name"

    // 函数调用
    fv := bean.CallFun("test", "Fun1", parameterMap)
    assert.Equal(t, isc.ToString(fv[0]), "name")
}
```

### 在线管理bean功能
在运行中如果出现问题，需要查看某个对象的属性和函数的时候，就可以使用该功能，进行动态的查看、修改对应属性，以及动态的执行对应的函数

```yaml
base:
  endpoint:
    # 是否启用bean的端点，默认false
    bean:
      enable: true
```

```shell
# 获取注册的所有bean
curl http://localhost:xxx/{api-prefix}/{api-module}/bean/name/all'
# 查询注册的某些bean 
curl http://localhost:xxx/{api-prefix}/{api-module}/bean/name/list/:name'
# 查询某个bean的属性值
curl -X POST http://localhost:xxx/{api-prefix}/{api-module}/bean/field/get' -d '{"bean": "xx", "field": "xxx"}'
# 修改某个bean的属性的值（暂时只支持基本类型）
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/bean/field/set' -d '{"bean": "xx", "field": "xxx", "value": "xxx"}'
# 调用bean的某个函数（参数暂时只支持基本类型）
curl -X POST http://localhost:xxx/{api-prefix}/{api-module}/bean/fun/call' -d '{"bean": "xx", "fun": "xxx", "parameter": {"p1":"xx", "p2": "xxx"}}'
```

提示：<br/>
- 调用bean函数中，parameter的对应的map中的key只能是p1、p2、p3...这种表示的是第一个、第二个、第三个参数的值
- 调用bean函数中，参数值暂时只适用于基本结构，对于实体类或者map类的暂时不支持，后续可以考虑支持


