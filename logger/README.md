# logger
The isc-gobase's logger package provides a fast and simple logger dedicated to format output.

logger API is designed to provide both a great developer experience and stunning [performance](#benchmarks). Its unique chaining API allows zerolog to write JSON (or CBOR) log events by avoiding allocations and reflection.

Uber's [zap](https://godoc.org/go.uber.org/zap) library pioneered this approach.

To keep the code base and the API simple, zerolog focuses on efficient structured logging only. Pretty logging on the console is made possible using the provided (but inefficient) [`zerolog.ConsoleWriter`](#pretty-logging).
## Installation
`go get github.com/isyscore/isc-gobase`
## Getting started
### simple Logging Example
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
```yaml
base:
  logger:
    # 日志root级别：trace/debug/info/warn/error/fatal/panic，默认：info
    level: info
    # 日志打印的文件路径：full：全路径，short：半路径。默认：半路径
    path: full/short
    time:
      # 时间格式，time包中的内容
      format: time.RFC3339
    # 日志颜色
    color:
      # 启用：true/false，默认：false
      enable: false
    split:
      # 日志是否启用切分：true/false，默认false
      enable: false
      # 日志拆分的单位：MB
      size: 300
    max:
      ## 日志文件最大保留天数
      history: 7
    ## 日志文件目录，默认工程目录的logs文件夹
    dir: ./logs/
    ## 是否将console信息打印到文件app-console.log，默认false
    console:
      writeFile: false

```

### 线上日志级别动态修改
支持线上动态的日志修改
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.logger.level", "value":"debug"}'
```

提示：<br/>
目前日志级别粒度比较粗，比如修改了级别为debug后，则大于等于debug的级别都会打印，粒度还是比较粗，建议后续增加日志分组概念
