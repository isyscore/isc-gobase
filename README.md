https://github.com/isyscore/isc-gobase/actions/workflows/go.yml/badge.svg

# isc-gobase

isc-gobase 框架是杭州指令集智能科技有限公司在java转go的实践中沉淀总结的一套至简化工具框架。遵从大道至简原则，让开发者在开发go的项目方面使用更简单

## 下载
```shell
go get github.com/isyscore/isc-gobase
```
提示：更新相关依赖
```shell
go mod tidy
```

## 快速入门
isc-gobase定位是工具框架，包含各种各样的工具，并对开发中的各种常用的方法进行封装。也包括web方面的工具
### web项目
创建`main.go`文件和同目录的`application.yml` 文件

```text
├── application.yaml
├── go.mod
└── main.go
```

```yaml
# application.yml 内容
base:
  server:
    # 是否启用，默认：false
    enable: true
```
```go
// main.go 文件
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/isyscore/isc-gobase/server"
    "github.com/isyscore/isc-gobase/server/rsp"
)

func main() {
    server.Get("api/get", GetData)
    server.Run()
}

func GetData(c *gin.Context) {
    rsp.SuccessOfStandard(c, "value")
}
```
运行如下
```shell
root@user ~> curl http://localhost:8080/api/get
{"code":0,"data":"value","message":"success"}
```

### 包列表
|包名        | 简介 |
| --------   | :----: |
| [isc](/isc)| 基础工具（更新中）|
| [config](/config)| 配置文件管理|
| [validate](/validate)|校验核查 |
| [logger](/logger)| 日志 |
| [database](/database)|数据库处理（待更新） |
| [server](/server)| 服务处理 |
| [goid](/goid)| 局部id传递处理（theadlocal） |
| [json](/json)| json字符串处理工具 |
| [cache](/cache)| 缓存工具 |
| [time](/time)| 时间管理工具 |
| [file](/file)| 文件管理工具 |
| [coder](/coder)| 编解码加解密工具 |
| [http](/http)| http的辅助工具 |
| [listener](/listener)| 事件监听机制 |
| [bean](/bean)| 对象管理工具 |
| [debug](/debug)| 线上调试工具统一介绍文档 |
| [extend/orm](/extend/orm)| gorm、xorm的封装 |
| [extend/etcd](/extend/etcd)| etcd封装 |
| [extend/redis](/extend/redis)| go-redis的封装 |

### isc-gobase 项目测试
根目录提供go_test.sh文件，统一执行所有gobase中包的测试模块
```shell
sh go_test.sh
```
