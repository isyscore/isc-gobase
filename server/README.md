
```yaml
server:
  exception:
    # 打印异常返回
    print:
      # 是否启用：true, false；默认 true
      enable: true
      # 一些异常code不打印；默认可不填
      except:
        - 408
        - 409
```
