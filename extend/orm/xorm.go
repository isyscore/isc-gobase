package orm

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/logger"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
	"xorm.io/xorm/log"
)

type GobaseXormHook interface {
	BeforeProcess(c *contexts.ContextHook, driverName string) (context.Context, error)
	AfterProcess(c *contexts.ContextHook, driverName string) error
}

var defaultXormHooks []DefaultXormHook

type DefaultXormHook struct {
	driverName string
	gobaseXormHook GobaseXormHook
}

func (defaultHook *DefaultXormHook)BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return defaultHook.gobaseXormHook.BeforeProcess(c, defaultHook.driverName)
}

func (defaultHook *DefaultXormHook)AfterProcess(c *contexts.ContextHook) error {
	return defaultHook.gobaseXormHook.AfterProcess(c, defaultHook.driverName)
}

func init() {
	defaultXormHooks = []DefaultXormHook{}
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

func AddXormHook(hook GobaseXormHook) {
	defaultXormHook := DefaultXormHook{gobaseXormHook: hook}
	defaultXormHooks = append(defaultXormHooks, defaultXormHook)
	xormDbs := bean.GetBeanWithNamePre(constants.BeanNameXormPre)
	if xormDbs == nil {
		return
	}
	for _, db := range xormDbs {
		db.(*xorm.Engine).AddHook(&defaultXormHook)
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

	for _, hook := range defaultXormHooks {
		hook.driverName = datasourceConfig.DriverName
		xormDb.AddHook(&hook)
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

	//xormDb.ShowSQL(true)
	//xormDb.SetLogger(&XormLoggerAdapter{})
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

// LoggerAdapter wraps a Logger interface as LoggerContext interface
type XormLoggerAdapter struct {
}

// BeforeSQL implements ContextLogger
func (l *XormLoggerAdapter) BeforeSQL(ctx log.LogContext) {}

// AfterSQL implements ContextLogger
func (l *XormLoggerAdapter) AfterSQL(ctx log.LogContext) {
	var sessionPart string
	v := ctx.Ctx.Value("__xorm_session_id")
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
	}
	if ctx.ExecuteTime > 0 {
		logger.Group("orm").Debug("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		logger.Group("orm").Debug("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	}
}

// Debugf implements ContextLogger
func (l *XormLoggerAdapter) Debugf(format string, v ...interface{}) {
	logger.Group("orm").Debug(format, v)
}

// Errorf implements ContextLogger
func (l *XormLoggerAdapter) Errorf(format string, v ...interface{}) {
	logger.Error(format, v)
}

// Infof implements ContextLogger
func (l *XormLoggerAdapter) Infof(format string, v ...interface{}) {
	logger.Info(format, v)
}

// Warnf implements ContextLogger
func (l *XormLoggerAdapter) Warnf(format string, v ...interface{}) {
	logger.Warn(format, v)
}

// Level implements ContextLogger
func (l *XormLoggerAdapter) Level() log.LogLevel {
	return log.LOG_INFO
}

// SetLevel implements ContextLogger
func (l *XormLoggerAdapter) SetLevel(lv log.LogLevel) {
}

// ShowSQL implements ContextLogger
func (l *XormLoggerAdapter) ShowSQL(show ...bool) {

}

// IsShowSQL implements ContextLogger
func (l *XormLoggerAdapter) IsShowSQL() bool {
	return true
}

