
# server
server包是用于更加方便的开发web项目而封装的包，开启配置的话，如下

### 简单示例
```yaml
# application.yml 内容
base:
  server:
    # 是否启用，默认：false
    enable: true
```
```go
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

## 全部配置

isc-gobase项目内置的一些server的配置
```yaml
api-module: sample

base:
  api:
    # api前缀，默认包含api前缀，如果路径本身有api，则不再添加api前缀
    prefix: /api
  server:
    # 是否启用，默认：true
    enable: true
    # 端口号，默认：8080
    port: 8080
    gin:
      # 有三种模式：debug/release/test，默认 release
      mode: debug
      pprof:
        # pprof开关是否可以开启，默认false
        enable: false
    request:
      print:
        # 是否打印：true, false；默认 false
        enable: false
        # 指定要打印请求的uri
        include-uri:
          - /xxx/x/xxx
          - /xxx/x/xxxy
        # 指定不打印请求的uri
        exclude-uri:
          - /xxx/x/xxx
          - /xxx/x/xxxy
    response:
      print:
        # 是否打印：true, false；默认 false
        enable: false
        # 指定要打印请求的uri
        include-uri:
          - /xxx/x/xxx
          - /xxx/x/xxxy
        # 指定不打印请求的uri
        exclude-uri:
          - /xxx/x/xxx
          - /xxx/x/xxxy
    exception:
      # 异常返回打印
      print: false
        # 是否启用：true, false；默认 false
        enable: true
        # 一些异常httpStatus不打印；默认可不填
        exclude:
          - 408
          - 409
    # 版本号设置,默认值:unknown
    version: 1.0.0
  swagger:
    # 是否开启swagger：true, false；默认 false
    enable: false
```
isc-gobase项目内置的一些endpoint端口
```shell
base:
  # 内部开放的 endpoint
  endpoint:
    # 健康检查处理，默认关闭，true/false
    health:
      enable: true
    # 配置的管理（查看和变更），默认关闭，true/false
    config:
      enable: true
    # bean的管理（属性查看、属性修改、函数调用），默认false
    bean:
      enable: true
```

### api.prefix和api-module介绍
其中api和api-module这个配置最后的url前缀是<br/>
{api.prefix}/{api-module}/业务代码

比如如上：

```shell
root@user ~> curl http://localhost:8080/api/sample/get/data
{"code":0,"data":"ok","message":"success"}
```

### server介绍
额外说明：
提供request和response的打印，用于调试时候使用
```shell
# 开启请求的打印，开启后默认打印所有请求，如果想打印指定uri，请先配置uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.enable", "value":"true"}'
# 开启响应的打印
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.enable", "value":"true"}'
# 开启异常的打印
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.exception.print.enable", "value":"true"}'
```

#### 指定uri打印
如果不指定uri则会默认打印所有的请求
```shell
# 指定要打印的请求的uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[0]", "value":"/api/xx/xxx"}'
# 指定不要打印的请求uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.exclude-uri[0]", "value":"/api/xx/xxx"}'

# 指定要打印的响应的uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[0]", "value":"/api/xx/xxx"}'
# 指定不要打印的响应uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.exclude-uri[0]", "value":"/api/xx/xxx"}'
```

提示：<br/>
- 如果"请求"和"响应"都开启打印，则只会打印"响应"，因为响应中已经包括了"请求"
- 指定多个uri的话，如下，配置其实是按照properties的方式进行指定的
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[0]", "value":"/api/xx/xxx"}'
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[1]", "value":"/api/xx/xxy"}'
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[2]", "value":"/api/xx/xxz"}'
...
```
