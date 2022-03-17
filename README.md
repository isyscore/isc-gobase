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
api-module: api/app/sample
server:
  # 是否启用，默认：true
  enable: false
  # 端口号
  port: 8080
  gin:
    # 有三种模式：debug/release/test
    mode: debug

base:
  application:
    # 应用名称
    name: sample
  endpoint:
    # 健康检查处理，默认关闭，true/false
    health:
      enable: true
    # 配置的动态实时变更，默认关闭，true/false
    config:
      enable: true
```
```go
// main.go 文件
import (
  "github.com/gin-gonic/gin"
  "github.com/isyscore/isc-gobase/server"
)

func main() {
    server.RegisterRoute("/api/app/demo/get/data", server.HmGet, func(c *gin.Context) {
        c.Data(200, "application/json; charset=utf-8", []byte("ok"))
    })

    // 简化版，自动添加api-model
    server.GetApiModel("/demo/get/data", func(c *gin.Context) {
        c.Data(200, "application/json; charset=utf-8", []byte("ok"))
    })
    server.Run()
}
```
### 非 web 项目
对于非web项目，isc-gobase是作为一个纯工具使用。可以不使用配置文件，也可以使用配置文件不启用server服务
```yaml
server:
  # 不启用server，默认：true
  enable: false
```

### 各包的用法
|包名        | 简介 |
| --------   | :----: |
| [isc](/isc)| 基础工具（待更新）|
| [config](/config)| 配置文件管理|
| [validate](/validate)|校验核查 |
| [logger](/logger)| 日志（待更新） |
| [coder](/coder)| 编解码（待更新） |
| [database](/database)|数据库处理（待更新） |
| [file](/file)| 文件处理（待更新） |
| [http](/http)| http处理（待更新） |
| [server](/server)| 服务处理（待更新） |

