package tracing

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/zipkin"
)
type GoBaseRedisHook struct {
	redis.Hook
}

func NewGoRedisTracer() redis.Hook {
	return GoBaseRedisHook{}
}

var spanKeyRedis = "gobase-redis-span"

func (GoBaseRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// 这里是关键，通过 envoy 传过来的 header 解析出父 span，如果没有，则会创建新的根 span
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	spanCtx, err := zipkinPropagator.Extract(opentracing.HTTPHeadersCarrier(GetHeader()))
	if err != nil {
		logger.Warn("span 解析失败, 错误原因: %v", err)
		return ctx, err
	}

	span, _ := opentracing.StartSpanFromContext(ctx, cmd.Name(), opentracing.ChildOf(spanCtx))
	ctx = context.WithValue(ctx, spanKeyRedis, span)
	return ctx, nil
}

func (GoBaseRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span, ok := ctx.Value(spanKeyRedis).(opentracing.Span)
	if !ok || span == nil {
		return nil
	}
	defer span.Finish()

	// 记录error
	err := cmd.Err()
	if err != nil {
		span.LogFields(opentracinglog.Error(err))
	}

	args, err := json.Marshal(cmd.Args())
	if err != nil {
		span.LogFields(opentracinglog.Error(err))
	}

	span.LogFields(
		opentracinglog.String("cmd", cmd.Name()),
		opentracinglog.String("fullName", cmd.FullName()),
		opentracinglog.String("parentId", GetHeaderWithKey("x-b3-spanid")),
		opentracinglog.String("parameters", string(args)),
	)
	return nil
}

