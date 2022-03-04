package config

var ServerCfg ServerConfig
var BaseCfg BaseConfig
var LogCfg LoggerConfig

// ServerConfig server前缀
type ServerConfig struct {
	Port   int
	Lookup bool
}

// BaseConfig base前缀
type BaseConfig struct {
	Application AppApplication
	Profiles    AppProfile
}

// LoggerConfig log 前缀
type LoggerConfig struct {
	Level string
	Path  string
}

type AppProfile struct {
	Active string
}

type AppApplication struct {
	Name string
}

type StorageConnectionConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Parameters string
}
