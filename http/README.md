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
