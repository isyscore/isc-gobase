# logger

```go
package main

import (
    "github.com/isyscore/isc-gobase/logger"
)

func main() {
    logger.Info("my test info %s","i am info")
    logger.Warn("my test warn %s","i am warn")
    logger.Error("my test error %s","i am error")
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
```

### 线上日志级别动态修改
支持线上动态的日志修改
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.logger.level", "value":"debug"}'
```

提示：<br/>
目前日志级别粒度比较粗，比如修改了级别为debug后，则大于等于debug的级别都会打印，粒度还是比较粗，建议后续增加日志分组概念
