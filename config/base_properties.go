package config

var ApiModule string
var BaseCfg BaseConfig

// BaseConfig base前缀
type BaseConfig struct {
	Api         BaseApi         `yaml:"api"`
	Application BaseApplication `yaml:"application"`
	Server      BaseServer      `yaml:"server"`
	EndPoint    BaseEndPoint    `yaml:"endpoint"`
	Logger      BaseLogger      `yaml:"logger"`
	Profiles    BaseProfile     `yaml:"profiles"`
}

type BaseApi struct {
	prefix string `yaml:"prefix"` // api前缀
}

type BaseApplication struct {
	Name string `yaml:"name"` // 应用名字
}

type BaseServer struct {
	Enable    bool          `yaml:"enable"`    // 是否启用
	Port      int           `yaml:"port"`      // 端口号
	Gin       BaseGin       `yaml:"gin"`       // web框架gin的配置
	Exception BaseException `yaml:"exception"` // 异常处理
}

type BaseGin struct {
	Mode string `yaml:"mode"` // 有三种模式：debug/release/test
}

type BaseEndPoint struct {
	Health EndPointHealth `yaml:"health"` // 健康检查[端点]
	Config EndPointConfig `yaml:"config"` // 配置管理[端点]
}

type EndPointHealth struct {
	Enable bool `yaml:"enable"` // 是否启用
}

type EndPointConfig struct {
	Enable bool `yaml:"enable"` // 是否启用
}

type BaseException struct {
	Print ExceptionPrint `yaml:"print"` // 异常返回打印
}

type ExceptionPrint struct {
	Enable bool  `yaml:"enable"` // 是否启用
	Except []int `yaml:"except"` // 排除的httpStatus；默认可不填
}

type BaseLogger struct {
	Level string      `yaml:"level"` // 日志root级别：trace/debug/info/warn/error/fatal/panic，默认：info
	Time  LoggerTime  `yaml:"time"`  // 时间配置
	Color LoggerColor `yaml:"color"` // 日志颜色
	Split LoggerSplit `yaml:"split"` // 日志切分
}

type LoggerTime struct {
	Format string `yaml:"format"` // 时间格式，time包中的内容，比如：time.RFC3339
}

type LoggerColor struct {
	Enable bool `yaml:"enable"` // 是否启用
}

type LoggerSplit struct {
	Enable bool `yaml:"enable"` // 日志是否启用切分：true/false，默认false
	Size   int  `yaml:"size"`   // 日志拆分的单位：MB
}

type BaseProfile struct {
	Active string `yaml:"active"`
}

type StorageConnectionConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Parameters string `yaml:"parameters"`
}
