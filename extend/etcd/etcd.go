package etcd

import (
	"context"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"time"
)

func init() {
	config.LoadConfig()

	if config.ExistConfigFile() && config.GetValueBoolDefault("base.etcd.enable", false) {
		err := config.GetValueObject("base.etcd", &config.EtcdCfg)
		if err != nil {
			logger.Warn("读取redis配置异常")
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

	// 添加埋点
	return &EtcdClientWrap{Client: etcdClient}, nil
	//return etcdClient, nil
}

func NewEtcdClientWithCfg(etcdCfg etcdClientV3.Config) (*etcdClientV3.Client, error) {
	etcdClient, err := etcdClientV3.New(etcdCfg)
	if err != nil {
		logger.Error("生成etcd-client失败：%v", err.Error())
		return nil, err
	}

	// 添加埋点
	if EtcdTracingIsOpen() {
		//
	}

	return etcdClient, nil
}

type EtcdClientWrap struct {
	*etcdClientV3.Client
}

func (etcdWrap *EtcdClientWrap) Put(ctx context.Context, key, val string, opts ...etcdClientV3.OpOption) (*etcdClientV3.PutResponse, error) {
	return etcdWrap.Client.Put(ctx, key, val, opts...)
}

func (etcdWrap *EtcdClientWrap) Get(ctx context.Context, key string, opts ...etcdClientV3.OpOption) (*etcdClientV3.GetResponse, error) {
	return etcdWrap.Client.Get(ctx, key, opts...)
}

func (etcdWrap *EtcdClientWrap) Delete(ctx context.Context, key string, opts ...etcdClientV3.OpOption) (*etcdClientV3.DeleteResponse, error) {
	return etcdWrap.Client.Delete(ctx, key, opts...)
}

func (etcdWrap *EtcdClientWrap) Compact(ctx context.Context, rev int64, opts ...etcdClientV3.CompactOption) (*etcdClientV3.CompactResponse, error) {
	return etcdWrap.Client.Compact(ctx, rev, opts...)
}

func (etcdWrap *EtcdClientWrap) Do(ctx context.Context, op etcdClientV3.Op) (etcdClientV3.OpResponse, error) {
	return etcdWrap.Client.Do(ctx, op)
}

func (etcdWrap *EtcdClientWrap) Txn(ctx context.Context) etcdClientV3.Txn {
	return etcdWrap.Client.Txn(ctx)
}

func (etcdWrap *EtcdClientWrap) MemberList(ctx context.Context) (*etcdClientV3.MemberListResponse, error) {
	return etcdWrap.Client.MemberList(ctx)
}

func (etcdWrap *EtcdClientWrap) MemberAdd(ctx context.Context, peerAddrs []string) (*etcdClientV3.MemberAddResponse, error) {
	return etcdWrap.Client.MemberAdd(ctx)
}

func (etcdWrap *EtcdClientWrap) MemberAddAsLearner(ctx context.Context, peerAddrs []string) (*etcdClientV3.MemberAddResponse, error) {
	return etcdWrap.Client.MemberAddAsLearner(ctx)
}

func (etcdWrap *EtcdClientWrap) MemberRemove(ctx context.Context, id uint64) (*etcdClientV3.MemberRemoveResponse, error) {
	return etcdWrap.Client.MemberRemove(ctx)
}

func (etcdWrap *EtcdClientWrap) MemberPromote(ctx context.Context, id uint64) (*etcdClientV3.MemberPromoteResponse, error) {
	return etcdWrap.Client.MemberPromote(ctx)
}

func (etcdWrap *EtcdClientWrap) Grant(ctx context.Context, ttl int64) (*etcdClientV3.LeaseGrantResponse, error) {
	return etcdWrap.Client.Grant(ctx)
}

func (etcdWrap **EtcdClientWrap) Revoke(ctx context.Context, id etcdClientV3.LeaseID) (*etcdClientV3.LeaseRevokeResponse, error) {
	return etcdWrap.Client.Revoke(ctx, id)
}

func (etcdWrap *EtcdClientWrap) TimeToLive(ctx context.Context, id etcdClientV3.LeaseID, opts ...etcdClientV3.LeaseOption) (*etcdClientV3.LeaseTimeToLiveResponse, error) {
	return etcdWrap.Client.TimeToLive(ctx)
}

func (etcdWrap *EtcdClientWrap) Leases(ctx context.Context) (*etcdClientV3.LeaseLeasesResponse, error) {
	return etcdWrap.Client.Leases(ctx)
}

func (etcdWrap *EtcdClientWrap) KeepAlive(ctx context.Context, id etcdClientV3.LeaseID) (<-chan *etcdClientV3.LeaseKeepAliveResponse, error) {
	return etcdWrap.Client.KeepAlive(ctx)
}

func (etcdWrap *EtcdClientWrap) KeepAliveOnce(ctx context.Context, id etcdClientV3.LeaseID) (*etcdClientV3.LeaseKeepAliveResponse, error) {
	return etcdWrap.Client.KeepAliveOnce(ctx, id)
}

func (etcdWrap *EtcdClientWrap) Close() error {
	return etcdWrap.Client.Close()
}

func (etcdWrap *EtcdClientWrap) Watch(ctx context.Context, key string, opts ...etcdClientV3.OpOption) etcdClientV3.WatchChan {
	return etcdWrap.Client.Watch()
}

func (etcdWrap *EtcdClientWrap) RequestProgress(ctx context.Context) error {
	return etcdWrap.Client.RequestProgress(ctx)
}

func (etcdWrap *EtcdClientWrap) Authenticate(ctx context.Context, name string, password string) (*etcdClientV3.AuthenticateResponse, error) {
	return etcdWrap.Client.Authenticate(ctx, name, password)
}

func (etcdWrap *EtcdClientWrap) AuthEnable(ctx context.Context) (*etcdClientV3.AuthenticateResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) AuthDisable(ctx context.Context) (*etcdClientV3.AuthDisableResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) AuthStatus(ctx context.Context) (*etcdClientV3.AuthStatusResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserAdd(ctx context.Context, name string, password string) (*etcdClientV3.AuthUserAddResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserAddWithOptions(ctx context.Context, name string, password string, opt *etcdClientV3.UserAddOptions) (*etcdClientV3.AuthUserAddResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserDelete(ctx context.Context, name string) (*etcdClientV3.AuthUserDeleteResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserChangePassword(ctx context.Context, name string, password string) (*etcdClientV3.AuthUserChangePasswordResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserGrantRole(ctx context.Context, user string, role string) (*etcdClientV3.AuthUserGrantRoleResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserGet(ctx context.Context, name string) (*etcdClientV3.AuthUserGetResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserList(ctx context.Context) (*etcdClientV3.AuthUserListResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) UserRevokeRole(ctx context.Context, name string, role string) (*etcdClientV3.AuthUserRevokeRoleResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) RoleAdd(ctx context.Context, name string) (*etcdClientV3.AuthRoleAddResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) RoleGrantPermission(ctx context.Context, name string, key, rangeEnd string, permType etcdClientV3.PermissionType) (*etcdClientV3.AuthRoleGrantPermissionResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) RoleGet(ctx context.Context, role string) (*etcdClientV3.AuthRoleGetResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) RoleList(ctx context.Context) (*etcdClientV3.AuthRoleListResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) RoleRevokePermission(ctx context.Context, role string, key, rangeEnd string) (*etcdClientV3.AuthRoleRevokePermissionResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) RoleDelete(ctx context.Context, role string) (*etcdClientV3.AuthRoleDeleteResponse, error) {
	return etcdWrap.Client.Authenticate()
}

func (etcdWrap *EtcdClientWrap) AlarmList(ctx context.Context) (*etcdClientV3.AlarmResponse, error) {

}

func (etcdWrap *EtcdClientWrap) AlarmDisarm(ctx context.Context, m *etcdClientV3.AlarmMember) (*etcdClientV3.AlarmResponse, error) {

}

func (etcdWrap *EtcdClientWrap) Defragment(ctx context.Context, endpoint string) (*etcdClientV3.DefragmentResponse, error) {

}

func (etcdWrap *EtcdClientWrap) Status(ctx context.Context, endpoint string) (*etcdClientV3.StatusResponse, error) {

}

func (etcdWrap *EtcdClientWrap) HashKV(ctx context.Context, endpoint string, rev int64) (*etcdClientV3.HashKVResponse, error) {

}

func (etcdWrap *EtcdClientWrap) Snapshot(ctx context.Context) (io.ReadCloser, error) {

}

func (etcdWrap *EtcdClientWrap) MoveLeader(ctx context.Context, transfereeID uint64) (*etcdClientV3.MoveLeaderResponse, error) {

}

func (etcdWrap *EtcdClientWrap) WithLogger(lg *zap.Logger) *etcdClientV3.Client {

}

func (etcdWrap *EtcdClientWrap) GetLogger() *zap.Logger {
	return etcdWrap.Client.GetLogger()
}

func (etcdWrap *EtcdClientWrap) Ctx() context.Context {
	return etcdWrap.Client.GetLogger()
}

func (etcdWrap *EtcdClientWrap) Endpoints() []string {
	return etcdWrap.Client.GetLogger()
}

func (etcdWrap *EtcdClientWrap) SetEndpoints(eps ...string) {
	return etcdWrap.Client.GetLogger()
}

func (etcdWrap *EtcdClientWrap) Sync(ctx context.Context) error {
	return etcdWrap.Client.GetLogger()
}

func (etcdWrap *EtcdClientWrap) Dial(ep string) (*grpc.ClientConn, error) {
	return etcdWrap.Client.GetLogger()
}

func (etcdWrap *EtcdClientWrap) ActiveConnection() *grpc.ClientConn {
	return etcdWrap.Client.GetLogger()
}

func EtcdTracingIsOpen() bool {
	return config.GetValueBoolDefault("base.tracing.enable", true) && config.GetValueBoolDefault("base.tracing.etcd.enable", false)
}
