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

```
