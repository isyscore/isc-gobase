# 版本

### v1.5.3（doing）
#### 新增
1. 新增：增加kafka的配置化

#### 优化


#### 修复


---

### v1.5.2
#### 新增
1. 新增：增加emqx的配置化
2. 新增：增加logger包的group的多值分组标签
3. 新增：增加config的配置转化对象中的兼容yaml和json的功能
4. 新增：增加orm的配置，sql动态打印，支持gorm和xorm

#### 优化
1. 优化：优化set和list方面的实现
2. 优化：增加数据库的配置，用于优化gorm默认不支持mariadb问题

#### 修复
1. 修复：修复gorm的启动不添加埋点下的异常
2. 修复：修复redis的默认值不生效问题

---

### v1.5.1
#### 新增
1. 新增：core跨域配置化处理
2. 新增：网络可达性api
3. 新增：http的配置化处理

---

### v1.5.0
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
9. 新增：validate包新增业务自定义参数传递功能
10. 新增：config包新增yaml、yml、json和properties配置文件的占位符功能
#### 优化
1. 调整：go-redis的包结构，不兼容旧版本
2. 优化：返回值code的处理
3. 优化：重构logger包，简化代码，优化功能，增加分组能力
#### 修复
1. 修复：goid多个实例获取同一个值的问题
2. 修复：config的配置变更中的复杂结构类型不支持问题
3. 修复：中英文的长度判断不一致问题
