# isc-gobase

isc-gobase 框架是杭州指令集智能科技有限公司在java转go的实践中沉淀总结的一套至简化工具框架。遵从大道至简原则，让开发者在开发go的项目方面使用更简单

## 下载
```shell
go get github.com/isyscore/isc-gobase
```

## 快速入门
isc-gobase定位是工具框架，包含各种各样的工具，并对开发中的各种常用的方法进行封装。也包括web方面的工具
### web项目
创建`main.go`文件和同目录的`application.yml` 文件

```yaml
# application.yml 内容
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
      mode: debug
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
    server.Get("group1/data", GetData)
    server.Run()
}

func GetData(c *gin.Context) {
    rsp.SuccessOfStandard(c, "ok")
}
```
运行如下
```shell
root@user ~> curl http://localhost:8080/api/app/sample/group1/data
{"code":"success","data":"ok","message":"成功"}
```
### 非 web 项目
对于非web项目，isc-gobase是作为一个纯工具使用。这时候要分情况
- 不使用config包，则不需要添加 application.yml 文件
- 使用config包，不启动server，则将base.server.enable关闭
```yaml
base:
  server:
    # 是否启用server，默认：true
    enable: false
```

### 各包的用法
|包名        | 简介 |
| --------   | :----: |
| [isc](/isc)| 基础工具（待更新）|
| [config](/config)| 配置文件管理|
| [validate](/validate)|校验核查 |
| [logger](/logger)| 日志 |
| [coder](/coder)| 编解码（待更新） |
| [database](/database)|数据库处理（待更新） |
| [file](/file)| 文件处理（待更新） |
| [http](/http)| http处理（待更新） |
| [server](/server)| 服务处理 |
| [goid](/goid)| 局部id传递处理 |

