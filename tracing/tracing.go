package tracing

import (
	"github.com/go-redis/redis/v8"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/logger"
	"gorm.io/gorm"
	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
)

var GormHooks []gorm.Plugin
var XormHooks []contexts.Hook
var RedisHooks []redis.Hook
var EtcdHooks []etcd.GobaseEtcdHook

func init() {
	GormHooks = []gorm.Plugin{}
	XormHooks = []contexts.Hook{}
	RedisHooks = []redis.Hook{}
	EtcdHooks = []etcd.GobaseEtcdHook{}
}

func AddGormHook(hook gorm.Plugin) {
	GormHooks = append(GormHooks, hook)
	gormDbs := bean.GetBeanWithNamePre(constants.BeanNameGormPre)
	for _, db := range gormDbs {
		gormDb := db.(*gorm.DB)
		err := gormDb.Use(hook)
		if err != nil {
			logger.Error("添加hook出错: %v", err.Error())
		}
	}
}

func AddXormHook(hook contexts.Hook) {
	XormHooks = append(XormHooks, hook)
	xormDbs := bean.GetBeanWithNamePre(constants.BeanNameXormPre)
	for _, db := range xormDbs {
		db.(*xorm.Engine).AddHook(hook)
	}
}

func AddRedisHook(hook redis.Hook) {
	RedisHooks = append(RedisHooks, hook)
	redisDb := bean.GetBeanWithNamePre(constants.BeanNameRedisPre)
	if len(redisDb) > 0 {
		rd := redisDb[0].(redis.UniversalClient)
		rd.AddHook(hook)
	}
}

func AddEtcdHook(hook etcd.GobaseEtcdHook) {
	EtcdHooks = append(EtcdHooks, hook)
	client := bean.GetBean(constants.BeanNameEtcdPre)
	etcdClient := client.(*etcd.EtcdClientWrap)
	etcdClient.AddHook(hook)
}
