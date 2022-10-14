
## 单数据源配置
```yaml
base:
  datasource:
    username: xxx
    password: xxx
    host: xxx
    port: xxx
    # 目前支持: mysql、postgresql、sqlite、sqlserver
    driver-name: xxx
    # sqlite的的数据库路径
    sql-path: xxx.db
    # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
    url-config:
      xxx: xxx
      yyy: yyy
    # 连接池配置
    connect-pool:
      # 最大空闲连接数
      max-idel-conns: 10
      # 最大连接数
      max-open-conns: 10
      # 连接可重用最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
      max-life-time: 10
      # 连接空闲的最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
      max-idle-time: 10
      
```

## 多数据源配置
```yml
base:
  datasource:
    name1:
      username: xxx
      password: xxx
      host: xxx
      port: xxx
      # 目前支持: mysql、postgresql、sqlite、sqlserver
      driver-name: xxx
      path: xxx.sql
      # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
      url-config:
        xxx: xxx
        yyy: yyy
      # 连接池配置
      connect-pool:
        # 最大空闲连接数
        max-idel-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
        max-life-time: 10
        # 连接空闲的最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
        max-idle-time: 10
    name2:
      username: xxx
      password: xxx
      host: xxx
      port: xxx
      # 目前支持: mysql、postgresql、sqlite、sqlserver
      driver-name: xxx
      path: xxx.sql
      # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
      url-config:
        xxx: xxx
        yyy: yyy
      # 连接池配置
      connect-pool:
        # 最大空闲连接数
        max-idel-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
        max-life-time: 10
        # 连接空闲的最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
        max-idle-time: 10
      
```
