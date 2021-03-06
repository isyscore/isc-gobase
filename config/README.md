## config包
config包主要用于加载和管理项目中配置文件中的内容，配置文件为"application"开头的格式

### 1. 配置文件路径
默认该文件与main函数所在的类同目录
```go
// 示例
- application.yml
- application-local.yml
```

### 2. 配置文件格式
支持yaml、yml、json、properties配置文件
优先级: json > properties > yaml > yml

### 3. 支持profile加载不同配置文件
格式：application-{profile}.yyy
其中profile对应的变量为：base.profiles.active
变量的设置可以有如下
- 本地配置
- 环境变量配置

优先级：环境变量 > 本地配置

![img.png](img.png)

#### 代码中读取指定环境配置
```go
// 配置环境
os.Setenv("base.profiles.active", "local")

// 然后再加载的时候就会加载local的配置文件
config.LoadConfig()
```

### 4. 内置的配置文件自动加载
目前内置的自动加载的配置文件有如下这些，后续随着工程越来越大会越来越多
```yaml
api-module: sample
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
      mode: debug
  endpoint:
    # 健康检查处理，默认关闭，true/false
    health:
      enable: true
    # 配置的动态实时变更，默认关闭，true/false
    config:
      enable: true
```


### 5. 支持直接获取配置值
config包中提供了各种类型的api，方便实时获取
```go
// 基本类型
config.getValueInt("xxx.xxx")
config.getValueInt32("xxx.xxx")
config.getValueInt64("xxx.xxx")
config.getValueBool("xxx.xxx")
config.getValueString("xxx.xxx")
// ...

// 结构类型
config.getValueObject("xxx.xxx", &xxx)
```
示例：
```go
var ServerCfg ServerConfig

// base前缀
type BaseConfig struct {
    Application AppApplication
    Data string
}

type AppApplication struct {
    Name string
}
```

```yaml
base:
  application:
    name: "xxx-local"
  data: "test"
```

```go
// 直接读取即可
config.getValueObject("base", &ServerCfg)
```

说明：
v1.0.12版本后，支持对配置的中划线支持，此外还支持更多配置
- 中划线：比如：data-base-user
- 小驼峰：比如：dataBaseUser
- 大驼峰：比如：DataBaseUser
- 下划线：比如：data_base_user

比如：
```yaml
key1:
  ok1:
    hao-de-ok: 12
    name-age: 32
  ok2:
    haoDeOk: 12
    nameAge: 32
  ok3:
    HaoDeOk: 12
    NameAge: 32
  ok4:
    hao_de_ok: 12
    name_age: 32
```
```go
type SmallEntity struct {
    HaoDeOk int
    NameAge int
}

// 可以读取到
func TestSmall(t *testing.T) {
    config.LoadConfig()

    entity1 := SmallEntity{}
    config.GetValueObject("key1.ok1", &entity1)
    assert.Equal(t, entity1.NameAge, 32)
    assert.Equal(t, entity1.HaoDeOk, 12)

    entity2 := SmallEntity{}
    config.GetValueObject("key1.ok2", &entity2)
    assert.Equal(t, entity2.NameAge, 32)
    assert.Equal(t, entity2.HaoDeOk, 12)

    entity3 := SmallEntity{}
    config.GetValueObject("key1.ok3", &entity3)
    assert.Equal(t, entity3.NameAge, 32)
    assert.Equal(t, entity3.HaoDeOk, 12)

    entity4 := SmallEntity{}
    config.GetValueObject("key1.ok4", &entity4)
    assert.Equal(t, entity4.NameAge, 32)
    assert.Equal(t, entity4.HaoDeOk, 12)
}
```

### 6. 支持文件的绝对和相对路径读取
配置路径默认是与main同目录，也支持绝对路径读取对应的配置，该api可以用于与运维同学约定的服务器路径位置
```go
// 相对路径
config.LoadConfigFromRelativePath(xx)

// 绝对路径
config.LoadConfigFromAbsPath(xx)
```

### 7. 支持配置的叠加，相对路径和绝对路径
在配置已经加载完毕后，需要对一些配置进行覆盖，比如运维这边有相关的需求时候
```go
// 相对路径
config.AppendConfigFromRelativePath(xx)

// 绝对路径
config.AppendConfigFromAbsPath(xx)
```

### 8. 支持自动读取cm文件
应用启动会默认读取/home/{base.application.name}/config/application-default.yml对应的内容并覆盖应用的配置中

### 9. 支持配置的在线查看以及实时变更

如下配置开启后，就可以在线查看应用的所有配置了
```yaml
base:
  endpoint:
    # 配置的动态实时变更，默认关闭
    config:
      enable: true/false
```

```shell
// 查看应用所有配置
curl http://localhost:xxx/{api-prefix}/{api-module}/config/values

// 查看应用的某个配置
curl http://localhost:xxx/{api-prefix}/{api-module}/config/value/{key}

// 修改应用的配置
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"xxx", "value":"yyyy"}'
```

提示：<br/>
修改应用的配置会发送配置变更事件"event_of_config_change"，如果想要对配置变更进行监听，请监听，示例：
```go
func xxxx() {
    // 添加配置变更事件的监听，listener.EventOfConfigChange是内置的"event_of_config_change"
    listener.AddListener(listener.EventOfConfigChange, ConfigChangeListener)
}

func ConfigChangeListener(event listener.BaseEvent) {
    ev := event.(listener.ConfigChangeEvent)
    if ev.Key == "xxx" {
        value := ev.Value
        // 你的配置变更处理代码
    }
}
```

---

#### 注意

- 配置实体化
  - 无法动态的变更
  - 不支持默认配置
- api实时调用
  - 配置可以动态的变更
  - 有默认的api
    
建议：配置使用时候建议使用config.GetXXXX()

其中动态变更只对api实时调用的方式有效

---
