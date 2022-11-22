package orm

import (
	"github.com/isyscore/isc-gobase/bean"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/constants"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/tracing"
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

	var gormDb *gorm.DB
	gormDb, err = gorm.Open(getDialect(datasourceConfig.DriverName, datasourceConfig), gormConfig)
	if err != nil {
		logger.Warn("获取数据库db异常：%v", err.Error())
		return nil, err
	}

	for _, hook := range tracing.GormHooks {
		err := gormDb.Use(hook)
		if err != nil {
			logger.Error("gorm添加hook出错: %v", err.Error())
		}
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
	bean.AddBean(constants.BeanNameGormPre + datasourceName, gormDb)
	return gormDb, nil
}

func getDialect(dbType string, datasourceConfig config.DatasourceConfig) gorm.Dialector {
	dsn := getDbDsn(dbType, datasourceConfig)
	switch dbType {
	case "mysql":
		return mysql.Open(dsn)
	case "postgresql":
		return postgres.Open(dsn)
	case "sqlite":
		return sqlite.Open(dsn)
	case "sqlserver":
		return sqlserver.Open(dsn)
	}
	return nil
}

// 特殊字符处理
func specialCharChange(url string) string {
	return strings.ReplaceAll(url, "/", "%2F")
}

func OrmTracingIsOpen() bool {
	return config.GetValueBoolDefault("base.tracing.enable", false) && config.GetValueBoolDefault("base.tracing.orm.enable", false)
}
