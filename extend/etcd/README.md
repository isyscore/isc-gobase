# etcd

对业内的etcd客户端进行配置化封装，用于简化获取

全部配置
```yaml
base:
  etcd:
    # etcd的服务ip:port列表   
    endpoints:
      - 10.10.10.10:2379
      - 10.10.10.11:2379
    # 用户
    username: userxxx
    # 密码
    password: xxxxx
    # 自动同步间隔：是用其最新成员更新端点的间隔；默认为0，即禁用自动同步；配置示例：1s、1000ms
    auto-sync-interval: 5s
    # 拨号超时：是指连接失败后的超时时间；配置示例：1s、1000ms
    dial-timeout: 5s
    # 拨号保持连接时间：是客户端ping服务器以查看传输是否连接的时间；配置示例：1s、1000ms
    dial-keep-alive-time: 5s
    # 拨号保持连接超时：是客户端等待响应保持连接探测的时间，如果在此时间内没有收到响应，则连接将被关闭；配置示例：1s、1000ms
    dial-keep-alive-timeout: 5s
    # 拨号重试策略: 默认为空：表示默认不重试；1、2、3...表示重试多少次；always：表示一直重试
    dial-retry: 1
    # 最大呼叫：发送MSG大小是客户端请求发送的字节限制；默认：(2MB)2 * 1024 * 1024
    max-call-send-msg-size: 2 * 1024 * 1024
    # 最大调用recv MSG大小是客户端响应接收限制；默认：math.MaxInt32
    max-call-recv-msg-size: 10000000
    # 当设置拒绝旧集群时，将拒绝在过时的集群上创建客户端
    reject-old-cluster: false
    # 设置允许无流时将允许客户端发送keepalive ping到服务器没有任何活动流rp cs
    permit-without-stream: false
```
提供etcd的client获取api `NewEtcdClient`
```go
func Test1(t *testing.T) {
    etcdClient, _ := etcd.NewEtcdClient()

    ctx := context.Background()
    etcdClient.Put(ctx, "gobase.k1", "testValue")
    rsp, _ := etcdClient.Get(ctx, "gobase.k1")
}
```

### 示例
配置示例
```yaml
base:
  etcd:
    # 是否启用etcd
    enable: true
    # etcd的服务ip:port列表
    endpoints:
      - 10.10.10.10:2379
    # 用户
    username: xxx
    # 密码
    password: xxx
    # 拨号超时：是指连接失败后的超时时间；配置示例：1s、1000ms
    dial-timeout: 5s
```
代码示例
```go
func Test1(t *testing.T) {
    etcdClient, _ := etcd.NewEtcdClient()
    
    ctx := context.Background()
    etcdClient.Put(ctx, "gobase.k1", "testValue")
    rsp, _ := etcdClient.Get(ctx, "gobase.k1")
}
```
