# isc-gobase


## config包
config包主要用于加载和管理项目中配置文件中的内容，配置文件为"application"开头的格式

#### 1. 配置文件路径
默认该文件与main函数所在的类同目录
```go
// 示例
- applicaiton.go
- application.yml
- application-local.yml
```

#### 2. 配置文件格式
支持yaml、yml、json、properties配置文件
优先级: json > properties > yaml > yml

#### 3. 支持profile加载不同配置文件
格式：application-{profile}.yyy
其中profile对应的变量为：base.profiles.active
变量的设置可以有如下
- 本地配置
- 启动参数配置
- 环境变量配置

优先级：本地配置 > 启动参数 > 环境变量


#### 5. 支持直接获取配置值
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

#### 6. 支持文件的绝对和相对路径读取
配置路径默认是与main同目录，也支持绝对路径读取对应的配置，该api可以用于与运维同学
```go
// 相对路径
config.LoadConfigFromRelativePath(xx)

// 绝对路径
config.LoadConfigFromAbsPath(xx)
```

#### 7. 支持配置的叠加，相对路径和绝对路径
在配置已经加载完毕后，需要对一些配置进行覆盖，比如运维这边有相关的需求时候
```go
// 相对路径
config.AppendConfigFromRelativePath(xx)

// 绝对路径
config.AppendConfigFromAbsPath(xx)
```
