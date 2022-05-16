## redis
对go-redis进行封装，用于简化配置使用

### 快速使用

```yaml
base:
  redis:
    enable: true
    standalone:
      addr: localhost:16379
```

```go
func TestConnect(t *testing.T) {
    // 直接获取即可 
    rdb, _ := redis.GetClient()
    
    ctx := context.Background()
    rdb.Set(ctx, "k1", "vv", time.Hour)
    rlt := rdb.Get(ctx, "k1")
    fmt.Println(rlt.Result())
}
```

redis的全部配置如下

```yaml
base:
  redis:
    # 是否启用redis，默认关闭
    enable: bool
    password: string
    username: string
    # 单节点模式
    standalone:
      addr: string # 数据库节点
      database: int
      network: string # 网络类型，tcp或者unix，默认tcp
      read-only: bool # 开启从节点的只读功能
    # （主从高可用）哨兵模式
    sentinel: 
      master: string # 哨兵的集群名字
      addrs: string,string # 哨兵节点地址
      database: int # 数据库节点
      sentinel-user: string # 哨兵用户
      sentinel-password: string # 哨兵密码
      slave-only: bool # 将所有命令路由到从属只读节点。
    # 集群模式
    cluster: 
      addrs: string,string # 节点地址
      max-redirects: int # 最大重定向次数
      read-only: bool # 开启从节点的只读功能
      route-by-latency: bool # 允许将只读命令路由到最近的主节点或从节点，它会自动启用 ReadOnly
      route-randomly: bool # 允许将只读命令路由到随机的主节点或从节点，它会自动启用 ReadOnly
    
    # 命令执行失败配置
    max-retries: int # 命令执行失败时候，最大重试次数，默认3次，-1（不是0）则不重试
    min-retry-backoff: int #（单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
    max-retry-backoff: int # （单位毫秒）命令执行失败时候，每次重试的最大回退时间，默认512毫秒，-1则禁止回退
    
    # 超时配置
    dial-timeout: int # （单位毫秒）超时：创建新链接的拨号超时时间，默认15秒
    read-timeout: int # （单位毫秒）超时：读超时，默认3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
    write-timeout: int # （单位毫秒）超时：写超时，默认是读超时3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒

    # 连接池相关配置
    pool-fifo: bool # 连接池类型：fifo：true;lifo：false;和lifo相比，fifo开销更高
    poll-size: int # 最大连接池大小：默认每个cpu核是10个连接，cpu核数可以根据函数runtime.GOMAXPROCS来配置，默认是runtime.NumCpu
    min-idleconns: int # 最小空闲连接数
    max-conn-age: int #（单位毫秒） 连接存活时长，默认不关闭
    pool-timeout: int #（单位毫秒）获取链接池中的链接都在忙，则等待对应的时间，默认读超时+1秒
    idle-timeout: int #（单位毫秒）空闲链接时间，超时则关闭，注意：该时间要小于服务端的超时时间，否则会出现拿到的链接失效问题，默认5分钟，-1表示禁用超时检查
    idle-check-frequency: int #（单位毫秒）空闲链接核查频率，默认1分钟。-1禁止空闲链接核查，即使配置了IdleTime也不行
```
