package tracing

import (
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type CustomCodecCallOption struct {
	grpc.EmptyCallOption
}

func (o CustomCodecCallOption) before(c *etcdClientV3.callInfo) error {
	c.codec = o.Codec
	return nil
}
func (o CustomCodecCallOption) after(c *callInfo, attempt *csAttempt) {}


func WaitForReady(waitForReady bool) grpc.CallOption {
	return FailFastCallOption{FailFast: !waitForReady}
}
