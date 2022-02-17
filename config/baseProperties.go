package config

var ServerCfg ServerConfig
var BaseCfg BaseConfig
var LogCfg LoggerConfig

// server前缀
type ServerConfig struct {
	Port   int
	Lookup bool
}

// base前缀
type BaseConfig struct {
	Application AppApplication
	Profiles    AppProfile
}

// log 前缀
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
