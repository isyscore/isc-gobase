# logger

## 基本用法
```go
package main

import (
    "github.com/isyscore/isc-gobase/logger"
)

func main() {
    logger.Debug("hello %v", "debug")
    logger.Info("hello %v", "info")
    logger.Warn("hello %v", "warn")
    logger.Error("hello %v", "error")
}
```
配置
```yaml
base:
  logger:
    level: info
    # 日志文件目录，默认工程目录的logs文件夹
    home: ./logs/
    color:
      # 启用：true/false，默认：false
      enable: false
    # 日志滚动策略
    rotate:
      # 日志滚动size；默认300MB
      max-size: 300MB
      # 日志文件最大保留天数；默认60天
      max-history: 60d
      # 多久滚动一次；默认一天
      time: 1d
    path:
      # 日志展示格式：full-全路径；short-短路径；默认short
      type: short
    # 指定分组
    group:
      demo1: info
      demo2: debug
```

## 更多用法
### 1. 分组打印
对日志进行精细化管控，这里增加日志分组功能，其中默认分组为"root"，我们可以指定我们自己的分组，在代码中不同的业务代码中使用不同的分组，这样在需要的情况下，我们可以对该模块的日志级别进行精细化控制
```go
package main

import (
    "github.com/isyscore/isc-gobase/logger"
)

func main() {
    logger.Group("group1").Info("hello", " ", "info")
    logger.Group("group1").Infof("hello %v", "info")
}
```
### 2. 线上日志级别动态修改
支持线上动态的日志修改，base.logger.level为默认分组的日志级别，如下修改为默认分组
#### 2.1 root默认分组修改
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.logger.level", "value":"debug"}'
```
```go
func main() {
    // 设置为debug后，如下的debug级别能够记录到日志中
    logger.Debug("hello %v", "debug1")
}
```

#### 2.2 指定分组修改
如下为指定group名字为xxx的，设置日志级别为debug
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.logger.group.xxxx", "value":"debug"}'
```
```go
func main() {
    // 如下的日志就会被打印出来，否则默认info级别是不打印的
    logger.Group("group1").Debugf("hello %v", "debug")
}
```
