package etcd

import (
	"context"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/tracing"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

func init() {
	config.LoadConfig()

	if config.ExistConfigFile() && config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			logger.Warn("读取etcd配置异常")
			return
		}
	}
}

func NewEtcdClient() (*EtcdClientWrap, error) {
	// 客户端配置
	etcdCfg := etcdClientV3.Config{
		Endpoints:           config.EtcdCfg.Endpoints,
		Username:            config.EtcdCfg.Username,
		Password:            config.EtcdCfg.Password,
		MaxCallSendMsgSize:  config.EtcdCfg.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:  config.EtcdCfg.MaxCallRecvMsgSize,
		RejectOldCluster:    config.EtcdCfg.RejectOldCluster,
		PermitWithoutStream: config.EtcdCfg.PermitWithoutStream,
	}

	if config.EtcdCfg.AutoSyncInterval != "" {
		t, err := time.ParseDuration(config.EtcdCfg.AutoSyncInterval)
		if err != nil {
			logger.Warn("读取配置【base.etcd.auto-sync-interval】异常", err)
		} else {
			etcdCfg.AutoSyncInterval = t
		}
	}

	if config.EtcdCfg.DialTimeout != "" {
		t, err := time.ParseDuration(config.EtcdCfg.DialTimeout)
		if err != nil {
			logger.Warn("读取配置【base.etcd.dial-timeout】异常", err)
		} else {
			etcdCfg.DialTimeout = t
		}
	}

	if config.EtcdCfg.DialKeepAliveTime != "" {
		t, err := time.ParseDuration(config.EtcdCfg.DialKeepAliveTime)
		if err != nil {
			logger.Warn("读取配置【base.etcd.dial-keep-alive-time】异常", err)
		} else {
			etcdCfg.DialKeepAliveTime = t
		}
	}

	if config.EtcdCfg.DialKeepAliveTimeout != "" {
		t, err := time.ParseDuration(config.EtcdCfg.DialKeepAliveTimeout)
		if err != nil {
			logger.Warn("读取配置【base.etcd.dial-keep-alive-timeout】异常", err)
		} else {
			etcdCfg.DialKeepAliveTimeout = t
		}
	}

	etcdClient, err := etcdClientV3.New(etcdCfg)
	if err != nil {
		logger.Error("生成etcd-client失败：%v", err.Error())
		return nil, err
	}

	if EtcdTracingIsOpen() {
		err := tracing.InitTracing()
		if err != nil {
			logger.Warn("链路全局初始化失败，go-redis不接入埋点，错误：%v", err.Error())
		} else {
			return &EtcdClientWrap{Client: etcdClient, etcdHook: &tracing.GobaseEtcdHook{}}, nil
		}
	}
	return &EtcdClientWrap{Client: etcdClient}, nil
}

func NewEtcdClientWithCfg(etcdCfg etcdClientV3.Config) (*EtcdClientWrap, error) {
	etcdClient, err := etcdClientV3.New(etcdCfg)
	if err != nil {
		logger.Error("生成etcd-client失败：%v", err.Error())
		return nil, err
	}

	if EtcdTracingIsOpen() {
		err := tracing.InitTracing()
		if err != nil {
			logger.Warn("链路全局初始化失败，go-redis不接入埋点，错误：%v", err.Error())
		} else {
			return &EtcdClientWrap{Client: etcdClient, etcdHook: &tracing.GobaseEtcdHook{}}, nil
		}
	}
	return &EtcdClientWrap{Client: etcdClient}, nil
}

type EtcdClientWrap struct {
	*etcdClientV3.Client
	etcdHook *tracing.GobaseEtcdHook
}

func (etcdWrap *EtcdClientWrap) Put(ctx context.Context, key, val string, opts ...etcdClientV3.OpOption) (*etcdClientV3.PutResponse, error) {
	if !EtcdTracingIsOpen() {
		return etcdWrap.Client.Put(ctx, key, val, opts...)
	}
	op := etcdClientV3.OpPut(key, val, opts...)
	ctx = etcdWrap.etcdHook.Before(ctx, op)
	rsp, err := etcdWrap.Client.Put(ctx, key, val, opts...)
	etcdWrap.etcdHook.After(ctx, op, rsp, err)
	return rsp, err
}

func (etcdWrap *EtcdClientWrap) Get(ctx context.Context, key string, opts ...etcdClientV3.OpOption) (*etcdClientV3.GetResponse, error) {
	if !EtcdTracingIsOpen() {
		return etcdWrap.Client.Get(ctx, key, opts...)
	}
	op := etcdClientV3.OpGet(key, opts...)
	etcdWrap.etcdHook.Before(ctx, op)
	rsp, err := etcdWrap.Client.Get(ctx, key, opts...)
	etcdWrap.etcdHook.After(ctx, op, rsp, err)
	return rsp, err
}

func (etcdWrap *EtcdClientWrap) Delete(ctx context.Context, key string, opts ...etcdClientV3.OpOption) (*etcdClientV3.DeleteResponse, error) {
	if !EtcdTracingIsOpen() {
		return etcdWrap.Client.Delete(ctx, key, opts...)
	}
	op := etcdClientV3.OpDelete(key, opts...)
	etcdWrap.etcdHook.Before(ctx, op)
	rsp, err := etcdWrap.Client.Delete(ctx, key, opts...)
	etcdWrap.etcdHook.After(ctx, op, rsp, err)
	return rsp, err
}

func (etcdWrap *EtcdClientWrap) Do(ctx context.Context, op etcdClientV3.Op) (etcdClientV3.OpResponse, error) {
	if !EtcdTracingIsOpen() {
		return etcdWrap.Client.Do(ctx, op)
	}
	etcdWrap.etcdHook.Before(ctx, op)
	rsp, err := etcdWrap.Client.Do(ctx, op)
	etcdWrap.etcdHook.After(ctx, op, rsp, err)
	return rsp, err
}

func EtcdTracingIsOpen() bool {
	return config.GetValueBoolDefault("base.tracing.enable", true) && config.GetValueBoolDefault("base.tracing.etcd.enable", false)
}
