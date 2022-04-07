
## server
server包是用于更加方便的开发web项目而封装的包，开启配置的话，如下

```go
// main.go 文件
import (
    "github.com/gin-gonic/gin"
    "github.com/isyscore/isc-gobase/server"
    "github.com/isyscore/isc-gobase/server/rsp"
)

func main() {
    server.Get("group1/data", GetData)
    server.Run()
}

```

```yaml
base:
  server:
    # 是否启用，默认：true
    enable: true
    # 端口号
    port: 8080
    # web框架gin的配置
    gin:
      # 有三种模式：debug/release/test
      mode: debug
    exception:
      # 异常返回打印
      print:
        # 是否启用：true, false；默认 true
        enable: true
        # 一些异常httpStatus不打印；默认可不填
        except:
          - 408
          - 409
    # 版本号设置,默认值:unknown
    version: 1.0.0
  # 内部开放的 endpoint
  endpoint:
    # 健康检查处理，默认关闭，true/false
    health:
      enable: true
    # 配置的管理（查看和变更），默认关闭，true/false
    config:
      enable: true
```


