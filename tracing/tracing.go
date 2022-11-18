package tracing

//
//var GormHooks []gorm.Plugin
//var XormHooks []contexts.Hook
//var RedisHooks []redis.Hook
//var EtcdHooks []etcd.GobaseEtcdHook
//
//func init() {
//	GormHooks = make([]gorm.Plugin, 1)
//	XormHooks = make([]contexts.Hook, 1)
//	RedisHooks = make([]redis.Hook, 1)
//	EtcdHooks = make([]etcd.GobaseEtcdHook, 1)
//}
//
//func AddGormHook(hook gorm.Plugin) {
//	GormHooks = append(GormHooks, hook)
//	if orm.OrmTracingIsOpen() {
//		gormDbs := bean.GetBeanWithNamePre(constants.BeanNameGormPre)
//		for _, db := range gormDbs {
//			gormDb := db.(*gorm.DB)
//			err := gormDb.Use(hook)
//			if err != nil {
//				logger.Error("添加hook出错: %v", err.Error())
//			}
//		}
//	}
//}
//
//func AddXormHook(hook contexts.Hook) {
//	XormHooks = append(XormHooks, hook)
//	if orm.OrmTracingIsOpen() {
//		xormDbs := bean.GetBeanWithNamePre(constants.BeanNameXormPre)
//		for _, db := range xormDbs {
//			db.(*xorm.Engine).AddHook(hook)
//		}
//	}
//}
//
//func AddRedisHook(hook redis.Hook) {
//	RedisHooks = append(RedisHooks, hook)
//	if gobaseRedis.RedisTracingIsOpen() {
//		parameterMap := map[string]any{}
//		parameterMap["p1"] = hook
//
//		bean.CallFun(constants.BeanNameRedisPre, "AddHook", parameterMap)
//	}
//}
//
//func AddEtcdHook(hook etcd.GobaseEtcdHook) {
//	EtcdHooks = append(EtcdHooks, hook)
//	if etcd.EtcdTracingIsOpen() {
//		client := bean.GetBean(constants.BeanNameEtcdPre)
//		etcdClient := client.(*etcd.EtcdClientWrap)
//		etcdClient.AddHook(hook)
//	}
//}
