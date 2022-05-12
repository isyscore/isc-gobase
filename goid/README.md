
## goid
goid包是go版的ThreadLocal，用于在多个协程之间传递数据

```go
// 变量分配
var TestIdLocal = goid.NewLocalStorage()

// 设置
TestIdLocal.Set(tenantId)

// 获取
TestIdLocal.Get()
```
注意：<br/>
在遇到协程的时候，请不要使用go的原生方式，请使用如下的方式，否则goid数据无法传递
```go
// 启用协程
goid.Go(func() { 
    // ...
})
```

