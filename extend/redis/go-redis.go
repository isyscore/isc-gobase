package redis

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/logger"
	baseTime "github.com/isyscore/isc-gobase/time"
	//"github.com/isyscore/isc-gobase/tracing"
	//"github.com/isyscore/isc-gobase/tracing"
	"time"
)

var RedisHooks []goredis.Hook

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
	RedisHooks = []goredis.Hook{}
}

func NewClient() (goredis.UniversalClient, error) {
	var rdbClient goredis.UniversalClient
	if config.RedisCfg.Sentinel.Master != "" {
		rdbClient = goredis.NewFailoverClient(getSentinelConfig())
	} else if len(config.RedisCfg.Cluster.Addrs) != 0 {
		rdbClient = goredis.NewClusterClient(getClusterConfig())
	} else {
		rdbClient = goredis.NewClient(getStandaloneConfig())
	}

	for _, hook := range RedisHooks {
		rdbClient.AddHook(hook)
	}
	bean.AddBean(constants.BeanNameRedisPre, &rdbClient)
	return rdbClient, nil
}

func AddRedisHook(hook goredis.Hook) {
	RedisHooks = append(RedisHooks, hook)
	redisDb := bean.GetBeanWithNamePre(constants.BeanNameRedisPre)
	if len(redisDb) > 0 {
		rd := redisDb[0].(goredis.UniversalClient)
		rd.AddHook(hook)
	}
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
