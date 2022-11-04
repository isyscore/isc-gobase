package tracing

import (
	"context"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/zipkin"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"reflect"
)

var spanKeyEtcd = "gobase-redis-span"

type GobaseEtcdHook struct {
}

func (pHook *GobaseEtcdHook) Before(ctx context.Context, op etcdClientV3.Op) context.Context {
	// 这里是关键，通过 envoy 传过来的 header 解析出父 span，如果没有，则会创建新的根 span
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	spanCtx, err := zipkinPropagator.Extract(opentracing.HTTPHeadersCarrier(GetHeader()))
	if err != nil {
		logger.Warn("span 解析失败, 错误原因: %v", err)
		return ctx
	}

	span, _ := opentracing.StartSpanFromContext(ctx, getCmd(op), opentracing.ChildOf(spanCtx))
	ctx = context.WithValue(ctx, spanKeyEtcd, span)
	return ctx
}

func (pHook *GobaseEtcdHook) After(ctx context.Context, op etcdClientV3.Op, pRsp any, err error) {
	span, ok := ctx.Value(spanKeyEtcd).(opentracing.Span)
	if !ok || span == nil {
		return
	}
	defer span.Finish()

	// 记录error
	if err != nil {
		span.LogFields(opentracinglog.Error(err))
	}

	logger.Info("发送redis的埋点信息")
	span.LogFields(
		opentracinglog.String("req", isc.ToJsonString(toRequestOp(op))),
		opentracinglog.String("rsp", isc.ToJsonString(pRsp)),
		opentracinglog.String("parentId", GetHeaderWithKey("x-b3-spanid")),
	)
	return
}

func toRequestOp(op etcdClientV3.Op) *pb.RequestOp {
	if op.IsGet() {
		return &pb.RequestOp{Request: &pb.RequestOp_RequestRange{RequestRange: toRangeRequest(op)}}
	} else if op.IsPut() {
		r := &pb.PutRequest{
			Key:    op.KeyBytes(),
			Value:  op.ValueBytes(),
			Lease:  int64(isc.GetPrivateFieldValue(reflect.ValueOf(&op), "leaseID").(etcdClientV3.LeaseID)),
			PrevKv: isc.GetPrivateFieldValue(reflect.ValueOf(&op), "prevKV").(bool),
		}
		return &pb.RequestOp{Request: &pb.RequestOp_RequestPut{RequestPut: r}}
	} else if op.IsDelete() {
		r := &pb.DeleteRangeRequest{
			Key:      op.KeyBytes(),
			RangeEnd: op.RangeBytes(),
			PrevKv:   isc.GetPrivateFieldValue(reflect.ValueOf(&op), "prevKV").(bool),
		}
		return &pb.RequestOp{Request: &pb.RequestOp_RequestDeleteRange{RequestDeleteRange: r}}
	}
	return nil
}

func toRangeRequest(op etcdClientV3.Op) *pb.RangeRequest {
	if !op.IsGet() {
		return nil
	}
	r := &pb.RangeRequest{
		Key:               op.KeyBytes(),
		RangeEnd:          op.RangeBytes(),
		Revision:          op.Rev(),
		Serializable:      op.IsSerializable(),
		KeysOnly:          op.IsKeysOnly(),
		CountOnly:         op.IsCountOnly(),
		MinModRevision:    op.MinModRev(),
		MaxModRevision:    op.MaxModRev(),
		MinCreateRevision: op.MinCreateRev(),
		MaxCreateRevision: op.MaxCreateRev(),
	}
	return r
}

func getCmd(op etcdClientV3.Op) string {
	if op.IsGet() {
		return "get"
	} else if op.IsPut() {
		return "put"
	} else if op.IsDelete() {
		return "delete"
	}
	return ""
}
