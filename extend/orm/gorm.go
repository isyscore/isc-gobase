package orm

import (
	"context"
	"database/sql"
	driverMysql "github.com/go-sql-driver/mysql"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/listener"
	baseLogger "github.com/isyscore/isc-gobase/logger"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

func NewGormDb() (*gorm.DB, error) {
	return doNewGormDb("", &gorm.Config{})
}

func NewGormDbWitConfig(gormConfig *gorm.Config) (*gorm.DB, error) {
	return doNewGormDb("", gormConfig)
}

func NewGormDbWithName(datasourceName string) (*gorm.DB, error) {
	return doNewGormDb(datasourceName, &gorm.Config{})
}

func NewGormDbWithNameAndConfig(datasourceName string, gormConfig *gorm.Config) (*gorm.DB, error) {
	return doNewGormDb(datasourceName, gormConfig)
}

func doNewGormDb(datasourceName string, gormConfig *gorm.Config) (*gorm.DB, error) {
	datasourceConfig := config.DatasourceConfig{}
	targetDatasourceName := "base.datasource"
	if datasourceName != "" {
		targetDatasourceName = "base.datasource." + datasourceName
	}
	err := config.GetValueObject(targetDatasourceName, &datasourceConfig)
	if err != nil {
		baseLogger.Warn("读取读取配置【datasource】异常")
		return nil, err
	}

	// 注册原生的sql的hook
	if len(gormHooks) != 0 {
		sqlRegister(datasourceConfig.DriverName)
	}

	var gormDb *gorm.DB
	dsn := getDbDsn(datasourceConfig.DriverName, datasourceConfig)
	gormDb, err = gorm.Open(getDialect(dsn, datasourceConfig.DriverName), gormConfig)
	if err != nil {
		baseLogger.Warn("获取数据库db异常：%v", err.Error())
		return nil, err
	}

	d, _ := gormDb.DB()

	maxIdleConns := config.GetValueInt("base.datasource.connect-pool.max-idle-conns")
	if maxIdleConns != 0 {
		// 设置空闲的最大连接数
		d.SetMaxIdleConns(maxIdleConns)
	}

	maxOpenConns := config.GetValueInt("base.datasource.connect-pool.max-open-conns")
	if maxOpenConns != 0 {
		// 设置数据库打开连接的最大数量
		d.SetMaxOpenConns(maxOpenConns)
	}

	maxLifeTime := config.GetValueString("base.datasource.connect-pool.max-life-time")
	if maxLifeTime != "" {
		// 设置连接可重复使用的最大时间
		t, err := time.ParseDuration(maxLifeTime)
		if err != nil {
			baseLogger.Warn("读取配置【base.datasource.connect-pool.max-life-time】异常", err)
		} else {
			d.SetConnMaxLifetime(t)
		}
	}

	maxIdleTime := config.GetValueString("base.datasource.connect-pool.max-idle-time")
	if maxIdleTime != "" {
		// 设置conn最大空闲时间设置连接空闲的最大时间
		t, err := time.ParseDuration(maxIdleTime)
		if err != nil {
			baseLogger.Warn("读取配置【base.datasource.connect-pool.max-idle-time】异常", err)
		} else {
			d.SetConnMaxIdleTime(t)
		}
	}

	gormDb.Logger = &GormLoggerAdapter{}
	bean.AddBean(constants.BeanNameGormPre+datasourceName, gormDb)
	// 添加orm的配置监听器
	listener.AddListener(listener.EventOfConfigChange, ConfigChangeListenerOfOrm)

	return gormDb, nil
}

// 特殊字符处理
func specialCharChange(url string) string {
	return strings.ReplaceAll(url, "/", "%2F")
}

func getDialect(dsn, driverName string) gorm.Dialector {
	switch driverName {
	case "mysql":
		return mysql.New(getMysqlConfig(dsn, driverName))
	case "postgresql":
		return postgres.New(postgres.Config{DSN: dsn, DriverName: WrapDriverName(driverName)})
	case "sqlite":
		return sqlite.Dialector{DSN: dsn, DriverName: WrapDriverName(driverName)}
	case "sqlserver":
		return sqlserver.New(sqlserver.Config{DSN: dsn, DriverName: WrapDriverName(driverName)})
	}
	return nil
}

func sqlRegister(driverName string) {
	name := WrapDriverName(driverName)
	for _, driver := range sql.Drivers() {
		if driver == name {
			return
		}
	}

	switch driverName {
	case "mysql":
		sql.Register(name, sqlhooks.Wrap(&driverMysql.MySQLDriver{}, &GobaseSqlHookProxy{DriverName: driverName}))
	case "postgresql":
		sql.Register(name, sqlhooks.Wrap(&pq.Driver{}, &GobaseSqlHookProxy{DriverName: driverName}))
	case "sqlite":
		sql.Register(name, sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &GobaseSqlHookProxy{DriverName: driverName}))
		//case "sqlserver": 暂时不支持
		//	sql.Register(WrapDriverName(driverName), sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &GobaseSqlHookProxy{}))
	}
}

func getMysqlConfig(dsn, driverName string) mysql.Config {
	return mysql.Config{
		DriverName: driverName,
		DSN: dsn,
		ServerVersion:                 config.GetValueStringDefault("base.datasource.mysql.server-version", ""),
		SkipInitializeWithVersion:     config.GetValueBoolDefault("base.datasource.mysql.skip-initialize-with-version", false),
		DefaultStringSize:             config.GetValueUIntDefault("base.datasource.mysql.default-string-size", 0),
		DisableWithReturning:          config.GetValueBoolDefault("base.datasource.mysql.disable-with-returning", false),
		DisableDatetimePrecision:      config.GetValueBoolDefault("base.datasource.mysql.disable-datetime-precision", false),
		DontSupportRenameIndex:        config.GetValueBoolDefault("base.datasource.mysql.dont-support-rename-index", false),
		DontSupportRenameColumn:       config.GetValueBoolDefault("base.datasource.mysql.dont-support-rename-column", false),
		DontSupportForShareClause:     config.GetValueBoolDefault("base.datasource.mysql.dont-support-for-share-clause", false),
		DontSupportNullAsDefaultValue: config.GetValueBoolDefault("base.datasource.mysql.dont-support-null-as-default-value", false),
	}
}

func WrapDriverName(driverName string) string {
	if len(gormHooks) != 0 {
		return driverName + "Hook"
	}
	return driverName
}

type GobaseGormHook interface {
	Before(ctx context.Context, driverName string, parameters map[string]any) (context.Context, error)
	After(ctx context.Context, driverName string, parameters map[string]any) (context.Context, error)
	Err(ctx context.Context, driverName string, err error, parameters map[string]any) error
}

var gormHooks []GobaseGormHook

func init() {
	gormHooks = []GobaseGormHook{}
}

func AddGormHook(hook GobaseGormHook) {
	gormHooks = append(gormHooks, hook)
}

type GobaseSqlHookProxy struct {
	DriverName string
}

func (proxy *GobaseSqlHookProxy) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	var ctxFinal context.Context
	for _, hook := range gormHooks {
		parametersMap := map[string]any{
			"query": query,
			"args":  args,
		}
		_ctx, err := hook.Before(ctx, proxy.DriverName, parametersMap)
		if err != nil {
			return _ctx, err
		} else {
			ctxFinal = _ctx
		}
	}
	return ctxFinal, nil
}

func (proxy *GobaseSqlHookProxy) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	for _, hook := range gormHooks {
		parametersMap := map[string]any{
			"query": query,
			"args":  args,
		}
		ctx, err := hook.After(ctx, proxy.DriverName, parametersMap)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func (proxy *GobaseSqlHookProxy) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	for _, hook := range gormHooks {
		parametersMap := map[string]any{
			"query": query,
			"args":  args,
		}
		err := hook.Err(ctx, proxy.DriverName, err, parametersMap)
		if err != nil {
			return err
		}
	}
	return nil
}

type GormLoggerAdapter struct {
}

func (l *GormLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	var levelStr logrus.Level
	switch level {
	case logger.Silent:
		levelStr = logrus.TraceLevel
	case logger.Error:
		levelStr = logrus.ErrorLevel
	case logger.Warn:
		levelStr = logrus.WarnLevel
	case logger.Info:
		levelStr = logrus.InfoLevel
	}
	baseLogger.Group("orm").SetLevel(levelStr)
	return l
}

func (l *GormLoggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	baseLogger.Info(msg, data)
}

func (l *GormLoggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	baseLogger.Warn(msg, data)
}

func (l *GormLoggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	baseLogger.Error(msg, data)
}

func (l *GormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sqlStr, rowsAffected := fc()
	if err != nil {
		baseLogger.Group("orm").Errorf("[SQL][%v]%s; error: %v", elapsed, sqlStr, err.Error())
	} else {
		baseLogger.Group("orm").Debugf("[SQL][%v][row:%v]%s", elapsed, rowsAffected, sqlStr)
	}
}
