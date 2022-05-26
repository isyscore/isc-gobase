## bean
对对象进行管理。

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



