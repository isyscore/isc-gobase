package redis

import (
	goredis "github.com/go-redis/redis/v8"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
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
	if config.RedisCfg.Sentinel.Master != "" {
		return goredis.NewFailoverClient(getSentinelConfig()), nil
	} else if len(config.RedisCfg.Cluster.Addrs) != 0 {
		return goredis.NewClusterClient(getClusterConfig()), nil
	} else {
		return goredis.NewClient(getStandaloneConfig()), nil
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
		Username: config.RedisCfg.Username,
		Password: config.RedisCfg.Password,

		MaxRetries: config.RedisCfg.MaxRetries,
		// todo
		//MinRetryBackoff: config.RedisCfg.MinRetryBackoff,
		//MaxRetryBackoff: config.RedisCfg.MaxRetryBackoff,

		//DialTimeout:  config.RedisCfg.DialTimeout,
		//ReadTimeout:  config.RedisCfg.ReadTimeout,
		//WriteTimeout: config.RedisCfg.WriteTimeout,

		PoolFIFO:     config.RedisCfg.PoolFIFO,
		PoolSize:     config.RedisCfg.PoolSize,
		MinIdleConns: config.RedisCfg.MinIdleConns,
		//MaxConnAge:         config.RedisCfg.MaxConnAge,
		//PoolTimeout:        config.RedisCfg.PoolTimeout,
		//IdleTimeout:        config.RedisCfg.IdleTimeout,
		//IdleCheckFrequency: config.RedisCfg.IdleCheckFrequency,
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

		MaxRetries: config.RedisCfg.MaxRetries,
		// todo
		//MinRetryBackoff: config.RedisCfg.MinRetryBackoff * time.Second,
		//MaxRetryBackoff: config.RedisCfg.MaxRetryBackoff,

		//DialTimeout:  config.RedisCfg.DialTimeout,
		//ReadTimeout:  config.RedisCfg.ReadTimeout,
		//WriteTimeout: config.RedisCfg.WriteTimeout,

		PoolFIFO:     config.RedisCfg.PoolFIFO,
		PoolSize:     config.RedisCfg.PoolSize,
		MinIdleConns: config.RedisCfg.MinIdleConns,
		//MaxConnAge:         config.RedisCfg.MaxConnAge,
		//PoolTimeout:        config.RedisCfg.PoolTimeout,
		//IdleTimeout:        config.RedisCfg.IdleTimeout,
		//IdleCheckFrequency: config.RedisCfg.IdleCheckFrequency,
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

		MaxRetries: config.RedisCfg.MaxRetries,
		// todo
		//MinRetryBackoff: config.RedisCfg.MinRetryBackoff,
		//MaxRetryBackoff: config.RedisCfg.MaxRetryBackoff,

		//DialTimeout:        config.RedisCfg.DialTimeout,
		//ReadTimeout:        config.RedisCfg.ReadTimeout,
		//WriteTimeout:       config.RedisCfg.WriteTimeout,
		PoolFIFO:     config.RedisCfg.PoolFIFO,
		PoolSize:     config.RedisCfg.PoolSize,
		MinIdleConns: config.RedisCfg.MinIdleConns,

		//MaxConnAge:         config.RedisCfg.MaxConnAge,
		//PoolTimeout:        config.RedisCfg.PoolTimeout,
		//IdleTimeout:        config.RedisCfg.IdleTimeout,
		//IdleCheckFrequency: config.RedisCfg.IdleCheckFrequency,
	}
	return redisConfig
}
