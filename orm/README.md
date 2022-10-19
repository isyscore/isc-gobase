# orm
对业内的常见Orm进行封装，进行方便使用，目前支持的有
- gorm
- xorm
- borm：基于database包进行的配置化封装

### 单数据源
#### 代码
```go
import "github.com/isyscore/isc-gobase/orm"

// gorm：获取默认配置库实例
orm.GetGormDb()

// gorm：获取默认配置库实例，自定义配置
orm.GetGormDbWitConfig(gormConfig *gorm.Config)

// xorm：获取默认配置库实例
orm.GetXormDb()

// borm：获取默认配置库实例
orm.GetBormDb()
```
#### 配置
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
      max-idle-conns: 10
      # 最大连接数
      max-open-conns: 10
      # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
      max-life-time: 10s
      # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
      max-idle-time: 10s
```

### 多数据源
#### 代码
```go
import "github.com/isyscore/isc-gobase/orm"

// gorm：根据数据源配置名获取库实例
orm.GetGormDbWithName(datasourceName string)

// gorm：根据数据源配置名获取库实例，自定义配置
orm.GetGormDbWithNameAndConfig(gormConfig *gorm.Config)

// xorm：根据数据源配置名获取库实例
orm.GetXormDbWithName(datasourceName string)

// borm：根据数据源配置名获取库实例
orm.GetBormDbWithName(datasourceName string)
```
#### 配置
```yml
base:
  datasource:
    # 数据源配置名1
    xxx-name1:
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
        max-idle-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
        max-life-time: 10s
        # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
        max-idle-time: 10s
    # 数据源配置名2
    xxx-name2:
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
        max-idle-conns: 10
        # 最大连接数
        max-open-conns: 10
        # 连接可重用最大时间；带字符（s：秒，m：分钟，h：小时）
        max-life-time: 10s
        # 连接空闲的最大时间；带字符（s：秒，m：分钟，h：小时）
        max-idle-time: 10s
```
### 示例：gorm
```go
func Test(t *testing.T) {
    db, _ := orm.GetGormDb()

    // 删除表
    db.Exec("drop table isc_demo.gobase_demo")

    //新增表
    db.Exec("CREATE TABLE gobase_demo(\n" +
    	"  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',\n" +
    	"  `name` char(20) NOT NULL COMMENT '名字',\n" +
    	"  `age` INT NOT NULL COMMENT '年龄',\n" +
    	"  `address` char(20) NOT NULL COMMENT '名字',\n" +
    	"  \n" +
    	"  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',\n" +
    	"  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',\n" +
    	"\n" +
    	"  PRIMARY KEY (`id`)\n" +
    	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='测试表'")

    // 新增数据
    db.Create(&GobaseDemo{Name: "zhou", Age: 18, Address: "杭州", CreateTime: time.Now(), UpdateTime: time.Now()})
    db.Create(&GobaseDemo{Name: "zhou", Age: 11, Address: "杭州2", CreateTime: time.Now(), UpdateTime: time.Now()})

    // 查询：一行
    var demo GobaseDemo
    db.First(&demo).Where("name=?", "zhou")
}
```

### 示例：xorm
// todo
```go

```
### 示例：borm
// todo
```go

```
