package config

var ApiModule string
var BaseCfg BaseConfig

// BaseConfig base前缀
type BaseConfig struct {
	// api
	Api BaseApi
	// application
	Application BaseApplication
	// server
	Server BaseServer
	// endpoint
	EndPoint BaseEndPoint
	// logger
	Logger BaseLogger
	// profiles
	Profiles BaseProfile
}

type BaseApi struct {
	// api前缀
	prefix string
}

type BaseApplication struct {
	// 应用名字
	Name string
}

type BaseServer struct {
	// 是否启用
	Enable bool
	// 端口号
	Port int
	// web框架gin的配置
	Gin BaseGin
	// 异常处理
	Exception BaseException
}

type BaseGin struct {
	// 有三种模式：debug/release/test
	Mode string
}

type BaseEndPoint struct {
	// 健康检查[端点]
	Health EndPointHealth
	// 配置管理[端点]
	Config EndPointConfig
}

type EndPointHealth struct {
	// 是否启用
	Enable bool
}

type EndPointConfig struct {
	// 是否启用
	Enable bool
}

type BaseException struct {
	// 异常返回打印
	Print ExceptionPrint
}

type ExceptionPrint struct {
	// 是否启用
	Enable bool
	// 排除的httpStatus；默认可不填
	Except []int
}

type BaseLogger struct {
	// 日志root级别：trace/debug/info/warn/error/fatal/panic，默认：info
	Level string
	// 时间配置
	Time LoggerTime
	// 日志颜色
	Color LoggerColor
	// 日志切分
	Split LoggerSplit
}

type LoggerTime struct {
	// 时间格式，time包中的内容，比如：time.RFC3339
	Format string
}

type LoggerColor struct {
	// 是否启用
	Enable bool
}

type LoggerSplit struct {
	// 日志是否启用切分：true/false，默认false
	Enable bool
	// 日志拆分的单位：MB
	Size int
}

type BaseProfile struct {
	Active string
}

type StorageConnectionConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Parameters string
}
