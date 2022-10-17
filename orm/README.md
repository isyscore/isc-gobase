
## 单数据源配置
```yaml
base:
  datasource:
    username: user
    password: passwd
    host: 10.33.33.33
    port: 8080
    # 目前支持: mysql、postgresql、sqlite、sqlserver
    driver-name: mysql
    # 数据库名
    db-name: xx_db
    # sqlite的的数据库路径
    sqlite-path: xxx.db
    # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
    url-config:
      xxx: xxxxx
      yyy: yyyyy
    # 连接池配置
    connect-pool:
      # 最大空闲连接数
      max-idel-conns: 10
      # 最大连接数
      max-open-conns: 10
      # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
      max-life-time: 10s
      # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
      max-idle-time: 10s
```

## 多数据源配置
```yml
base:
  datasource:
    # 数据源配置名1
    name1:
      username: xxx
      password: xxx
      host: xxx
      port: xxx
      # 目前支持: mysql、postgresql、sqlite、sqlserver
      driver-name: xxx
      # 数据库名
      db-name: xx_db
      # sqlite的的数据库路径
      sqlite-path: xxx.db
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
        # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
        max-life-time: 10s
        # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
        max-idle-time: 10s
    # 数据源配置名2
    name2:
      username: xxx
      password: xxx
      host: xxx
      port: xxx
      # 目前支持: mysql、postgresql、sqlite、sqlserver
      driver-name: xxx
      # 数据库名
      db-name: xx_db
      # sqlite的的数据库路径
      sqlite-path: xxx.db
      # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
      url-config:
        xxx: xxxxx
        yyy: yyyyy
      # 连接池配置
      connect-pool:
        # 最大空闲连接数
        max-idel-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
        max-life-time: 10s
        # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
        max-idle-time: 10s
      
```
