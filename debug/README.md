## 线上调试文档介绍
鉴于提供的一些调试工具比较分散，这里进行统一介绍。<br/>



### 一、日志动态修改
在出现问题时候，希望开启debug日志，这里提供了动态修改的功能，如下
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.logger.level", "value":"debug"}'
```

提示：<br/>
目前日志级别粒度比较粗，比如修改了级别为debug后，则大于等于debug的级别都会打印，粒度还是比较粗，建议后续增加日志分组概念


### 二、接口动态打印
在出现问题时候如何确定接口是否正常，要开启请求和响应的话，就方便多了，这里提供了该功能

#### 指定uri
如果不指定uri则会默认打印所有的请求
```shell
# 指定要打印的请求的uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[0]", "value":"/api/xx/xxx"}'
# 指定不要打印的请求uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.exclude-uri[0]", "value":"/api/xx/xxx"}'

# 指定要打印的响应的uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[0]", "value":"/api/xx/xxx"}'
# 指定不要打印的响应uri
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.exclude-uri[0]", "value":"/api/xx/xxx"}'
```

提示：<br/>
- 如果"请求"和"响应"都开启打印，则只会打印"响应"，因为响应中已经包括了"请求"
- 指定多个uri的话，如下，配置其实是按照properties的方式进行指定的
```shell
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[0]", "value":"/api/xx/xxx"}'
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[1]", "value":"/api/xx/xxy"}'
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.include-uri[2]", "value":"/api/xx/xxz"}'
...
```
#### 开启打印
注意：开启打印的话，所有的接口都会打印，所以建议请先指定uri
```shell
# 开启请求的打印
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.enable", "value":"true"}'
# 开启响应的打印
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.request.print.enable", "value":"true"}'
# 开启异常的打印
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/config/update -d '{"key":"base.server.exception.print.enable", "value":"true"}'
```
### 三、动态管理bean
在运行中如果出现问题，需要查看某个对象的属性和函数的时候，就可以使用该功能，进行动态的查看、修改对应属性，以及动态的执行对应的函数

```shell
# 获取注册的所有bean
curl http://localhost:xxx/{api-prefix}/{api-module}/bean/name/all'
# 查询注册的某些bean 
curl http://localhost:xxx/{api-prefix}/{api-module}/bean/name/list/:name'
# 查询某个bean的属性值
curl -X POST http://localhost:xxx/{api-prefix}/{api-module}/bean/field/get' -d '{"bean": "xx", "field": "xxx"}'
# 修改某个bean的属性的值
curl -X PUT http://localhost:xxx/{api-prefix}/{api-module}/bean/field/set' -d '{"bean": "xx", "field": "xxx", "value": "xxx"}'
# 调用bean的某个函数
curl -X POST http://localhost:xxx/{api-prefix}/{api-module}/bean/fun/call' -d '{"bean": "xx", "fun": "xxx", "parameter": {"p1":"xx", "p2": "xxx"}}'
```

注意：<br/>
- 调用bean函数中，parameter的对应的map中的key只能是p1、p2、p3...这种表示的是第一个、第二个、第三个参数的值
- 调用bean函数中，参数值暂时只适用于基本结构，对于实体类或者map类的暂时不支持，后续可以考虑支持
