
## server
server包是用于更加方便的开发web项目而封装的包，开启配置的话，如下

```go

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
```
