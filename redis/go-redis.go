package redis

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	baseTime "github.com/isyscore/isc-gobase/time"
	"github.com/isyscore/isc-gobase/tracing"
	"time"
)

type ConfigError struct {
	ErrMsg string
}

func (error *ConfigError) Error() string {
	return error.ErrMsg
}

func init() {
	config.LoadConfig()

	if config.ExistConfigFile() && config.GetValueBoolDefault("base.redis.enable", false) {
		err := config.GetValueObject("base.redis", &config.RedisCfg)
		if err != nil {
			logger.Warn("读取redis配置异常")
			return
		}
	}
}

func GetClient() (goredis.UniversalClient, error) {
	var rdbClient goredis.UniversalClient
	if config.RedisCfg.Sentinel.Master != "" {
		rdbClient = goredis.NewFailoverClient(getSentinelConfig())
	} else if len(config.RedisCfg.Cluster.Addrs) != 0 {
		rdbClient = goredis.NewClusterClient(getClusterConfig())
	} else {
		rdbClient = goredis.NewClient(getStandaloneConfig())
	}

	if RedisTracingIsOpen() {
		err := tracing.InitTracing()
		if err != nil {
			logger.Warn("链路全局初始化失败，go-redis不接入埋点，错误：%v", err.Error())
		} else {
			rdbClient.AddHook(tracing.NewGoRedisTracer())
		}
	}
	return rdbClient, nil
}

func RedisTracingIsOpen() bool {
	return config.GetValueBoolDefault("base.tracing.enable", true) && config.GetValueBoolDefault("base.tracing.redis.enable", false)
}

func getStandaloneConfig() *goredis.Options {
	addr := "127.0.0.1:6379"
	if config.RedisCfg.Standalone.Addr != "" {
		addr = config.RedisCfg.Standalone.Addr
	}

	redisConfig := &goredis.Options{
		Addr: addr,

		DB:       config.RedisCfg.Standalone.Database,
		Network:  config.RedisCfg.Standalone.Network,
		Username: config.RedisCfg.Username,
		Password: config.RedisCfg.Password,

		MaxRetries:      config.RedisCfg.MaxRetries,
		MinRetryBackoff: baseTime.NumToTimeDuration(config.RedisCfg.MinRetryBackoff, time.Millisecond),
		MaxRetryBackoff: baseTime.NumToTimeDuration(config.RedisCfg.MaxRetryBackoff, time.Millisecond),

		DialTimeout:  baseTime.NumToTimeDuration(config.RedisCfg.DialTimeout, time.Millisecond),
		ReadTimeout:  baseTime.NumToTimeDuration(config.RedisCfg.ReadTimeout, time.Millisecond),
		WriteTimeout: baseTime.NumToTimeDuration(config.RedisCfg.WriteTimeout, time.Millisecond),

		PoolFIFO:           config.RedisCfg.PoolFIFO,
		PoolSize:           config.RedisCfg.PoolSize,
		MinIdleConns:       config.RedisCfg.MinIdleConns,
		MaxConnAge:         baseTime.NumToTimeDuration(config.RedisCfg.MaxConnAge, time.Millisecond),
		PoolTimeout:        baseTime.NumToTimeDuration(config.RedisCfg.PoolTimeout, time.Millisecond),
		IdleTimeout:        baseTime.NumToTimeDuration(config.RedisCfg.IdleTimeout, time.Millisecond),
		IdleCheckFrequency: baseTime.NumToTimeDuration(config.RedisCfg.IdleCheckFrequency, time.Millisecond),
	}
	return redisConfig
}

func getSentinelConfig() *goredis.FailoverOptions {
	redisConfig := &goredis.FailoverOptions{
		SentinelAddrs: config.RedisCfg.Sentinel.Addrs,
		MasterName:    config.RedisCfg.Sentinel.Master,

		DB:               config.RedisCfg.Sentinel.Database,
		Username:         config.RedisCfg.Username,
		Password:         config.RedisCfg.Password,
		SentinelUsername: config.RedisCfg.Sentinel.SentinelUser,
		SentinelPassword: config.RedisCfg.Sentinel.SentinelPassword,

		MaxRetries:      config.RedisCfg.MaxRetries,
		MinRetryBackoff: baseTime.NumToTimeDuration(config.RedisCfg.MinRetryBackoff, time.Millisecond),
		MaxRetryBackoff: baseTime.NumToTimeDuration(config.RedisCfg.MaxRetryBackoff, time.Millisecond),

		DialTimeout:  baseTime.NumToTimeDuration(config.RedisCfg.DialTimeout, time.Millisecond),
		ReadTimeout:  baseTime.NumToTimeDuration(config.RedisCfg.ReadTimeout, time.Millisecond),
		WriteTimeout: baseTime.NumToTimeDuration(config.RedisCfg.WriteTimeout, time.Millisecond),

		PoolFIFO:           config.RedisCfg.PoolFIFO,
		PoolSize:           config.RedisCfg.PoolSize,
		MinIdleConns:       config.RedisCfg.MinIdleConns,
		MaxConnAge:         baseTime.NumToTimeDuration(config.RedisCfg.MaxConnAge, time.Millisecond),
		PoolTimeout:        baseTime.NumToTimeDuration(config.RedisCfg.PoolTimeout, time.Millisecond),
		IdleTimeout:        baseTime.NumToTimeDuration(config.RedisCfg.IdleTimeout, time.Millisecond),
		IdleCheckFrequency: baseTime.NumToTimeDuration(config.RedisCfg.IdleCheckFrequency, time.Millisecond),
	}

	return redisConfig
}

func getClusterConfig() *goredis.ClusterOptions {
	if len(config.RedisCfg.Cluster.Addrs) == 0 {
		config.RedisCfg.Cluster.Addrs = []string{"127.0.0.1:6379"}
	}

	redisConfig := &goredis.ClusterOptions{
		Addrs: config.RedisCfg.Cluster.Addrs,

		Username: config.RedisCfg.Username,
		Password: config.RedisCfg.Password,

		MaxRedirects:   config.RedisCfg.Cluster.MaxRedirects,
		ReadOnly:       config.RedisCfg.Cluster.ReadOnly,
		RouteByLatency: config.RedisCfg.Cluster.RouteByLatency,
		RouteRandomly:  config.RedisCfg.Cluster.RouteRandomly,

		MaxRetries:      config.RedisCfg.MaxRetries,
		MinRetryBackoff: baseTime.NumToTimeDuration(config.RedisCfg.MinRetryBackoff, time.Millisecond),
		MaxRetryBackoff: baseTime.NumToTimeDuration(config.RedisCfg.MaxRetryBackoff, time.Millisecond),

		DialTimeout:  baseTime.NumToTimeDuration(config.RedisCfg.DialTimeout, time.Millisecond),
		ReadTimeout:  baseTime.NumToTimeDuration(config.RedisCfg.ReadTimeout, time.Millisecond),
		WriteTimeout: baseTime.NumToTimeDuration(config.RedisCfg.WriteTimeout, time.Millisecond),
		PoolFIFO:     config.RedisCfg.PoolFIFO,
		PoolSize:     config.RedisCfg.PoolSize,
		MinIdleConns: config.RedisCfg.MinIdleConns,

		MaxConnAge:         baseTime.NumToTimeDuration(config.RedisCfg.MaxConnAge, time.Millisecond),
		PoolTimeout:        baseTime.NumToTimeDuration(config.RedisCfg.PoolTimeout, time.Millisecond),
		IdleTimeout:        baseTime.NumToTimeDuration(config.RedisCfg.IdleTimeout, time.Millisecond),
		IdleCheckFrequency: baseTime.NumToTimeDuration(config.RedisCfg.IdleCheckFrequency, time.Millisecond),
	}
	return redisConfig
}
