package orm

import (
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/logger"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
)

var XormHooks []contexts.Hook

func init() {
	XormHooks = []contexts.Hook{}
}

func NewXormDb() (*xorm.Engine, error) {
	return doNewXormDb("", map[string]string{})
}

func NewXormDbWithParams(params map[string]string) (*xorm.Engine, error) {
	return doNewXormDb("", params)
}

func NewXormDbWithName(datasourceName string) (*xorm.Engine, error) {
	return doNewXormDb(datasourceName, map[string]string{})
}

func NewXormDbWithNameParams(datasourceName string, params map[string]string) (*xorm.Engine, error) {
	return doNewXormDb(datasourceName, params)
}

func AddXormHook(hook contexts.Hook) {
	XormHooks = append(XormHooks, hook)
	xormDbs := bean.GetBeanWithNamePre(constants.BeanNameXormPre)
	if xormDbs == nil {
		return
	}
	for _, db := range xormDbs {
		db.(*xorm.Engine).AddHook(hook)
	}
}

func doNewXormDb(datasourceName string, params map[string]string) (*xorm.Engine, error) {
	datasourceConfig := config.DatasourceConfig{}
	targetDatasourceName := "base.datasource"
	if datasourceName != "" {
		targetDatasourceName = "base.datasource." + datasourceName
	}
	err := config.GetValueObject(targetDatasourceName, &datasourceConfig)
	if err != nil {
		logger.Warn("读取读取配置【datasource】异常")
		return nil, err
	}

	var dsn = getDbDsn(datasourceConfig.DriverName, datasourceConfig)
	var xormDb *xorm.Engine
	xormDb, err = xorm.NewEngineWithParams(datasourceConfig.DriverName, dsn, params)
	if err != nil {
		logger.Warn("获取数据库db异常：%v", err.Error())
		return nil, err
	}

	for _, hook := range XormHooks {
		xormDb.AddHook(hook)
	}

	maxIdleConns := config.GetValueInt("base.datasource.connect-pool.max-idle-conns")
	if maxIdleConns != 0 {
		// 设置空闲的最大连接数
		xormDb.SetMaxIdleConns(maxIdleConns)
	}

	maxOpenConns := config.GetValueInt("base.datasource.connect-pool.max-open-conns")
	if maxOpenConns != 0 {
		// 设置数据库打开连接的最大数量
		xormDb.SetMaxOpenConns(maxOpenConns)
	}

	maxLifeTime := config.GetValueString("base.datasource.connect-pool.max-life-time")
	if maxLifeTime != "" {
		// 设置连接可重复使用的最大时间
		t, err := time.ParseDuration(maxLifeTime)
		if err != nil {
			logger.Warn("读取配置【base.datasource.connect-pool.max-life-time】异常", err)
		} else {
			xormDb.SetConnMaxLifetime(t)
		}
	}
	bean.AddBean(constants.BeanNameXormPre + datasourceName, xormDb)
	return xormDb, nil
}

func NewXormDbMasterSlave(masterDatasourceName string, slaveDatasourceNames []string, policies ...xorm.GroupPolicy) (*xorm.EngineGroup, error) {
	masterDb, err := NewXormDbWithName(masterDatasourceName)
	if err != nil {
		logger.Warn("获取数据库 主节点【%v】失败，%v", masterDatasourceName, err.Error())
		return nil, err
	}

	var slaveDbs []*xorm.Engine
	for _, slaveDatasource := range slaveDatasourceNames {
		slaveDb, err := NewXormDbWithName(slaveDatasource)
		if err != nil {
			logger.Warn("获取数据库 从节点【%v】失败，%v", slaveDatasource, err.Error())
			return nil, err
		}

		slaveDbs = append(slaveDbs, slaveDb)
	}

	return xorm.NewEngineGroup(masterDb, slaveDbs, policies...)
}
