# json

json包是为了方便对json格式的字符处理使用

```json
{
    "k1":12,
    "k2":true,
    "k3":{
        "k31":32,
        "k32":"str",
        "k33":{
            "k331":12
        }
    }
}
```

```go
func TestLoad(t *testing.T) {
    jsonObject := json.Object{}
    // str就是上面的字符串
    err := jsonObject.Load(str)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    assert.Equal(t, jsonObject.Get("k1"), 12)
    assert.Equal(t, jsonObject.Get("k2"), true)
    // 支持连key获取
    assert.Equal(t, jsonObject.Get("k3.k31"), 32)
    assert.Equal(t, jsonObject.Get("k3.k32"), "str")
    assert.Equal(t, jsonObject.Get("k3.k33.k331"), 12)
}
```

另外也提供了各种类型的api用于指定的类型使用
```go

GetString(key string) string
GetInt(key string) int
GetInt8(key string) int8
GetInt16(key string) int16
GetInt32(key string) int32
GetInt64(key string) int64
GetUInt(key string) uint
GetUInt8(key string) uint8
GetUInt16(key string) uint16
GetUInt32(key string) uint32
GetUInt64(key string) uint64
GetFloat32(key string) float32
GetFloat64(key string) float64
GetBool(key string) bool
GetObject(key string, targetPtrObj any) error
GetArray(key string) []any
```
