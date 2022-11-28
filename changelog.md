# 版本

### 分支：feature/trace
#### 新增
1. 新增：orm包的封装：xorm、gorm的配置化封装
2. 新增：etcd包的封装：go-etcd的配置化封装
3. 新增：debug的全局调试功能
3. 新增：listener包的分组功能
4. 新增：head的透传数据上下文
5. 新增：gorm、xorm、go-etcd、go-redis包的执行钩子回调功能
6. 新增：统一的异常回调功能
7. 新增：logger包的日志记录格式中的traceId和userId
8. 新增：server的最前缀filter功能
#### 优化
1. 调整：go-redis的包结构，不兼容
2. 优化：返回值code的处理
3. 优化：返回值code的处理
4. 调整：logger的日志
#### 修复
1. 修复：goid多个实例获取同一个值的问题
2. 修复config的配置变更中的复杂结构类型不支持问题


### 分支：feature/v1.4.4
#### 新增
1. 新增：validate包新增业务自定义参数传递功能
#### 修复
1. 修复：中英文的长度判断不一致问题
