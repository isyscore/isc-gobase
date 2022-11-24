package orm

import (
	"context"
	"database/sql"
	driverMysql "github.com/go-sql-driver/mysql"
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	"github.com/qustavo/sqlhooks/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
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
		logger.Warn("读取读取配置【datasource】异常")
		return nil, err
	}

	// 注册原生的sql的hook
	sqlRegister(datasourceConfig.DriverName)

	var gormDb *gorm.DB
	dsn := getDbDsn(datasourceConfig.DriverName, datasourceConfig)
	gormDb, err = gorm.Open(getDialect(dsn, datasourceConfig.DriverName), gormConfig)
	if err != nil {
		logger.Warn("获取数据库db异常：%v", err.Error())
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
			logger.Warn("读取配置【base.datasource.connect-pool.max-life-time】异常", err)
		} else {
			d.SetConnMaxLifetime(t)
		}
	}

	maxIdleTime := config.GetValueString("base.datasource.connect-pool.max-idle-time")
	if maxIdleTime != "" {
		// 设置conn最大空闲时间设置连接空闲的最大时间
		t, err := time.ParseDuration(maxIdleTime)
		if err != nil {
			logger.Warn("读取配置【base.datasource.connect-pool.max-idle-time】异常", err)
		} else {
			d.SetConnMaxIdleTime(t)
		}
	}
	bean.AddBean(constants.BeanNameGormPre+datasourceName, gormDb)
	return gormDb, nil
}

// 特殊字符处理
func specialCharChange(url string) string {
	return strings.ReplaceAll(url, "/", "%2F")
}


func getDialect(dsn, driverName string) gorm.Dialector {
	switch driverName {
	case "mysql":
		return mysql.New(mysql.Config{DSN: dsn, DriverName: WrapDriverName(driverName)})
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
	switch driverName {
	case "mysql":
		sql.Register(WrapDriverName(driverName), sqlhooks.Wrap(&driverMysql.MySQLDriver{}, &GobaseSqlHookProxy{}))
	case "postgresql":
		sql.Register(WrapDriverName(driverName), sqlhooks.Wrap(&pq.Driver{}, &GobaseSqlHookProxy{}))
	case "sqlite":
		sql.Register(WrapDriverName(driverName), sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &GobaseSqlHookProxy{}))
		//case "sqlserver": 暂时不支持
		//	sql.Register(WrapDriverName(driverName), sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, &GobaseSqlHookProxy{}))
	}
}

func WrapDriverName(driverName string) string {
	return driverName + "Hook"
}

type GobaseGormHook interface {
	Before(ctx context.Context, parameters map[string]any) (context.Context, error)
	After(ctx context.Context, parameters map[string]any) (context.Context, error)
	Err(ctx context.Context, err error, parameters map[string]any) error
}

var gormHooks []GobaseGormHook

func init() {
	gormHooks = []GobaseGormHook{}
}

func AddGormHook(hook GobaseGormHook) {
	gormHooks = append(gormHooks, hook)
}

type GobaseSqlHookProxy struct {}

func (proxy *GobaseSqlHookProxy) Before(ctx context.Context, query string, args ...interface {}) (context.Context, error) {
	for _, hook := range gormHooks {
		parametersMap := map[string]any{
			"query":query,
			"args": args,
		}
		ctx, err := hook.Before(ctx, parametersMap)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func (proxy *GobaseSqlHookProxy) After(ctx context.Context, query string, args ...interface {}) (context.Context, error) {
	for _, hook := range gormHooks {
		parametersMap := map[string]any{
			"query":query,
			"args": args,
		}
		ctx, err := hook.After(ctx, parametersMap)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

func (proxy *GobaseSqlHookProxy) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	for _, hook := range gormHooks {
		parametersMap := map[string]any{
			"query":query,
			"args": args,
		}
		err := hook.Err(ctx, err, parametersMap)
		if err != nil {
			return err
		}
	}
	return nil
}
