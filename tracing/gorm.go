package tracing

import (
	"context"
	"encoding/json"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go/zipkin"
	"gorm.io/gorm"
)

type GobaseGormHook struct {
}

const (
	spanKeyGorm = "gobase-gorm-span"

	// 自定义事件名称
	_eventBeforeCreate = "gobase-gorm-collector-event:before_create"
	_eventAfterCreate  = "gobase-gorm-collector-event:after_create"
	_eventBeforeUpdate = "gobase-gorm-collector-event:before_update"
	_eventAfterUpdate  = "gobase-gorm-collector-event:after_update"
	_eventBeforeQuery  = "gobase-gorm-collector-event:before_query"
	_eventAfterQuery   = "gobase-gorm-collector-event:after_query"
	_eventBeforeDelete = "gobase-gorm-collector-event:before_delete"
	_eventAfterDelete  = "gobase-gorm-collector-event:after_delete"
	_eventBeforeRow    = "gobase-gorm-collector-event:before_row"
	_eventAfterRow     = "gobase-gorm-collector-event:after_row"
	_eventBeforeRaw    = "gobase-gorm-collector-event:before_raw"
	_eventAfterRaw     = "gobase-gorm-collector-event:after_raw"

	// 自定义 span 的操作名称
	_opCreate = "insert"
	_opUpdate = "update"
	_opQuery  = "select"
	_opDelete = "delete"
	_opRow    = "row"
	_opRaw    = "execute"
)

// 开箱即用，serviceName: 此项目的微服务名称，collectorEndpoint: 数据收集器的地址(如:http://isc-core-back-service:31300/api/core/back/v1/middle/spans)
func NewGormPlugin() gorm.Plugin {
	return &GobaseGormHook{}
}

// 实现 gorm 插件所需方法
func (i *GobaseGormHook) Name() string {
	return "gobase_gorm_plugin"
}

// 实现 gorm 插件所需方法
func (i *GobaseGormHook) Initialize(db *gorm.DB) (err error) {
	// 在 gorm 中注册各种回调事件
	for _, e := range []error{
		db.Callback().Create().Before("gorm:create").Register(_eventBeforeCreate, beforeCreate),
		db.Callback().Create().After("gorm:create").Register(_eventAfterCreate, after),
		db.Callback().Update().Before("gorm:update").Register(_eventBeforeUpdate, beforeUpdate),
		db.Callback().Update().After("gorm:update").Register(_eventAfterUpdate, after),
		db.Callback().Query().Before("gorm:query").Register(_eventBeforeQuery, beforeQuery),
		db.Callback().Query().After("gorm:query").Register(_eventAfterQuery, after),
		db.Callback().Delete().Before("gorm:delete").Register(_eventBeforeDelete, beforeDelete),
		db.Callback().Delete().After("gorm:delete").Register(_eventAfterDelete, after),
		db.Callback().Row().Before("gorm:row").Register(_eventBeforeRow, beforeRow),
		db.Callback().Row().After("gorm:row").Register(_eventAfterRow, after),
		db.Callback().Raw().Before("gorm:raw").Register(_eventBeforeRaw, beforeRaw),
		db.Callback().Raw().After("gorm:raw").Register(_eventAfterRaw, after),
	} {
		if e != nil {
			return e
		}
	}
	return
}

// 注册各种前置事件时，对应的事件方法
func _injectBefore(db *gorm.DB, op string) {
	if db == nil {
		return
	}

	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "未定义 db.Statement 或 db.Statement.Context")
		return
	}

	// 这里是关键，通过 envoy 传过来的 header 解析出父 span，如果没有，则会创建新的根 span
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	spanCtx, err := zipkinPropagator.Extract(opentracing.HTTPHeadersCarrier(GetHeader()))
	if err != nil {
		logger.Error("jaeger span 解析失败, 错误原因: %v", err)
		return
	}
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, op, opentracing.ChildOf(spanCtx))
	db.InstanceSet(spanKeyGorm, span)
}

// 注册后置事件时，对应的事件方法
func after(db *gorm.DB) {
	if db == nil {
		return
	}

	if db.Statement == nil || db.Statement.Context == nil {
		db.Logger.Error(context.TODO(), "未定义 db.Statement 或 db.Statement.Context")
		return
	}

	_span, isExist := db.InstanceGet(spanKeyGorm)
	if !isExist || _span == nil {
		return
	}

	// 断言，进行类型转换
	span, ok := _span.(opentracing.Span)
	if !ok || span == nil {
		return
	}
	defer span.Finish()

	// 记录error
	if db.Error != nil {
		span.LogFields(opentracinglog.Error(db.Error))
	}

	b, err := json.Marshal(db.Statement.Vars)
	if err != nil {
		span.LogFields(opentracinglog.Error(err))
	}

	// 记录其他内容
	span.LogFields(
		opentracinglog.String("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)),
		opentracinglog.String("table", db.Statement.Table),
		opentracinglog.String("query", db.Statement.SQL.String()),
		opentracinglog.String("parentSpanId", GetHeaderWithKey("x-b3-spanid")),
		opentracinglog.String("parameters", string(b)),
	)
}

func beforeCreate(db *gorm.DB) {
	_injectBefore(db, _opCreate)
}

func beforeUpdate(db *gorm.DB) {
	_injectBefore(db, _opUpdate)
}

func beforeQuery(db *gorm.DB) {
	_injectBefore(db, _opQuery)
}

func beforeDelete(db *gorm.DB) {
	_injectBefore(db, _opDelete)
}

func beforeRow(db *gorm.DB) {
	_injectBefore(db, _opRow)
}

func beforeRaw(db *gorm.DB) {
	_injectBefore(db, _opRaw)
}
