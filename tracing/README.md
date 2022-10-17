# trace

### gorm
```yaml
base:
  tracing:
    # 链路中间件采集总开关；默认开启；
    enable: true
    # 服务搜集库；默认开启；
    collector-endpoint: http://isc-core-back-service:31300/api/core/back/v1/middle/spans
    orm:
      # 是否启动gorm采集；该开关与总开关是与的关系；目前go只支持gorm和xorm；默认关闭
      enable: false
```
