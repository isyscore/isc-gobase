## http
http包是简化封装解析，在公司这边有标准的返回结构，但是我们很多时候只使用其中的data部分，很多时候自己需要判断并解析，这里对这个做了简单的工具封装
```json
{
  "code": 0,
  "message": "yy",
  "data": "xx"
}
```

### 示例
```go
func TestGetSimple(t *testing.T) {
    // {"code":"success","data":"ok","message":"成功"}
    data, err := http.GetSimple("http://localhost:8082/api/api/app/sample/test/get")
    if err != nil {
        fmt.Errorf(err.Error())
        return
    }
    fmt.Println(string(data))

    // "ok"
    data, err = http.GetSimpleOfStandard("http://localhost:8082/api/api/app/sample/test/get")
    if err != nil {
        fmt.Errorf(err.Error())
        return
    }
    fmt.Println(string(data))
}
```

### 配置
```yaml
// todo
# http的配置
base:
  http:
    # 通用的超时配置，链接、重定向、响应的所有超时时间；默认0，就是不超时
    timeout: 0s
    transport:
      # 握手的最长时间
      tls-handshake-timeout: 1s
      # 禁用HTTP keep alives，且将仅对单个HTTP请求使用到服务器的连接。 这与类似命名的TCP keep-alives无关
      disable-keep-alives: true
      # 禁用gzip的压缩标示
      disable-compression: true
      # 最大空闲连接数
      max-idle-conns: 0
      # 每个主机最大空闲连接数
      max-idle-conns-per-host: 0
      # 每个主机最大连接数
      max-conns-per-host: 12
      # 连接在关闭之前保持空闲的最长时间
      idle-conn-timeout: 12s
      # 完全写入请求后等待服务器响应标头的时间
      response-header-timeout: 12s
      # 在请求具有“Expect:100 continue”标头时，在完全写入请求标头后等待服务器的第一个响应标头的时间
      expect-continue-timeout: 12s
      # 指定服务器响应标头中允许的响应字节数限制
      max-response-header-bytes: 123
      # 写入缓冲区的大小；如果为零，则使用默认值（当前为4KB）
      write-buffer-size: 0
      # 从传输读取时使用的读取缓冲区的大小。如果为零，则使用默认值（当前为4KB）
      read-buffer-size: 0
      # 使用Dial、DialTLS或DialContext func或TLSClientConfig字段时候，默认关闭http2；如果想要开启，则请设置为true
      ForceAttemptHTTP2: true
      # 用于创建未加密TCP连接
      dial-context:
        # 超时是拨号等待连接完成的最长时间。如果同时设置了Deadline，则可能会更早失败。 默认值为无超时。
        timeout: 0s
        # 超时的绝对时间
        deadline: "2023-02-13"
        # 活动网络连接的保持活动探测之间的间隔；默认15s
        keep-alive: 10s
```
