package config

var ApiModule string
var BaseCfg BaseConfig
var RedisCfg RedisConfig

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
	Prefix string `yaml:"prefix"` // api前缀
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

// ---------------------------- redis ----------------------------
// base.redis前缀
type RedisConfig struct {
	Password string
	Username string

	// 单节点
	Standalone RedisStandaloneConfig
	// 哨兵
	Sentinel RedisSentinelConfig
	// 集群
	Cluster RedisClusterConfig

	// ----- 命令执行失败配置 -----
	// 命令执行失败时候，最大重试次数，默认3次，-1（不是0）则不重试
	MaxRetries int
	// （单位毫秒） 命令执行失败时候，每次重试的最小回退时间，默认8毫秒，-1则禁止回退
	MinRetryBackoff int
	// （单位毫秒）命令执行失败时候，每次重试的最大回退时间，默认512毫秒，-1则禁止回退
	MaxRetryBackoff int

	// ----- 超时配置 -----
	// （单位毫秒）超时：创建新链接的拨号超时时间，默认15秒
	DialTimeout int
	// （单位毫秒）超时：读超时，默认3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
	ReadTimeout int
	// （单位毫秒）超时：写超时，默认是读超时3秒，使用-1，使用-1则表示无超时，0的话是表示默认3秒
	WriteTimeout int

	// ----- 连接池相关配置 -----
	// 连接池类型：fifo：true;lifo：false;和lifo相比，fifo开销更高
	PoolFIFO bool
	// 最大连接池大小：默认每个cpu核是10个连接，cpu核数可以根据函数runtime.GOMAXPROCS来配置，默认是runtime.NumCpu
	PoolSize int
	// 最小空闲连接数
	MinIdleConns int
	// （单位毫秒） 连接存活时长，默认不关闭
	MaxConnAge int
	// （单位毫秒）获取链接池中的链接都在忙，则等待对应的时间，默认读超时+1秒
	PoolTimeout int
	// （单位毫秒）空闲链接时间，超时则关闭，注意：该时间要小于服务端的超时时间，否则会出现拿到的链接失效问题，默认5分钟，-1表示禁用超时检查
	IdleTimeout int
	// （单位毫秒）空闲链接核查频率，默认1分钟。-1禁止空闲链接核查，即使配置了IdleTime也不行
	IdleCheckFrequency int
}

// base.redis.standalone
type RedisStandaloneConfig struct {
	Addr     string
	Database int
	// 网络类型，tcp或者unix，默认tcp
	Network  string `match:"value={tcp, unix}"  errMsg:"network值不合法，只可为两个值：tcp和unix"`
	ReadOnly bool
}

// base.redis.sentinel
type RedisSentinelConfig struct {
	// 哨兵的集群名字
	Master string
	// 哨兵节点地址
	Addrs []string
	// 数据库节点
	Database int
	// 哨兵用户
	SentinelUser string
	// 哨兵密码
	SentinelPassword string
	// 将所有命令路由到从属只读节点。
	SlaveOnly bool
}

type RedisClusterConfig struct {
	// 节点地址
	Addrs []string
	// 最大重定向次数
	MaxRedirects int
	// 开启从节点的只读功能
	ReadOnly bool
	// 允许将只读命令路由到最近的主节点或从节点，它会自动启用 ReadOnly
	RouteByLatency bool
	// 允许将只读命令路由到随机的主节点或从节点，它会自动启用 ReadOnly
	RouteRandomly bool
}
