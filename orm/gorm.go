package orm

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
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

func GetGormDb() (*gorm.DB, error) {
	return doGetGormDb("", &gorm.Config{})
}

func GetGormDbWitConfig(gormConfig *gorm.Config) (*gorm.DB, error) {
	return doGetGormDb("", gormConfig)
}

func GetGormDbWithName(datasourceName string) (*gorm.DB, error) {
	return doGetGormDb(datasourceName, &gorm.Config{})
}

func GetGormDbWithNameAndConfig(datasourceName string, gormConfig *gorm.Config) (*gorm.DB, error) {
	return doGetGormDb(datasourceName, gormConfig)
}

func doGetGormDb(datasourceName string, gormConfig *gorm.Config) (*gorm.DB, error) {
	datasourceConfig := DatasourceConfig{}
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

	if OrmTracingIsOpen() {
		err := tracing.InitTracing()
		if err != nil {
			logger.Warn("链路全局初始化失败，gorm 不接入埋点，错误：%v", err.Error())
		} else {
			err := gormDb.Use(tracing.NewGormPlugin())
			if err != nil {
				logger.Warn("接入tracing异常：%v", err.Error())
			}
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
		}
		d.SetConnMaxLifetime(t)
	}

	maxIdleTime := config.GetValueString("base.datasource.connect-pool.max-idle-time")
	if maxIdleTime != "" {
		// 设置conn最大空闲时间设置连接空闲的最大时间
		t, err := time.ParseDuration(maxIdleTime)
		if err != nil {
			logger.Warn("读取配置【base.datasource.connect-pool.max-idle-time】异常", err)
		}
		d.SetConnMaxIdleTime(t)
	}
	return gormDb, nil
}

func getDialect(dbType string, datasourceConfig DatasourceConfig) gorm.Dialector {
	sqlConfigMap := map[string]string{}
	err := config.GetValueObject("base.datasource.url-config", &sqlConfigMap)
	if err != nil {
		logger.Warn("读取配置【base.datasource.url-config】异常", err)
	}

	switch dbType {
	case "mysql":
		// 格式：user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.Host, datasourceConfig.Port, datasourceConfig.DbName)
		if len(sqlConfigMap) != 0 {
			var kvList []string
			for key, value := range sqlConfigMap {
				kvList = append(kvList, fmt.Sprintf("%s=%s", key, specialCharChange(value)))
			}
			dsn += fmt.Sprintf("?%s", strings.Join(kvList, "&"))
		}
		return mysql.Open(dsn)
	case "postgresql":
		// 格式：host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", datasourceConfig.Host, datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.DbName, datasourceConfig.Port)
		if len(sqlConfigMap) != 0 {
			var kvList []string
			for key, value := range sqlConfigMap {
				kvList = append(kvList, fmt.Sprintf("%s=%s", key, value))
			}
			dsn += fmt.Sprintf(" %s", strings.Join(kvList, " "))
		}
		return postgres.Open(dsn)
	case "sqlite":
		// 格式： gorm.db
		return sqlite.Open(datasourceConfig.SqlitePath)
	case "sqlserver":
		// 格式：sqlserver://user:password@localhost:9930?database=gorm
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.Host, datasourceConfig.Port, datasourceConfig.DbName)
		return sqlserver.Open(dsn)
	}
	return nil
}

// 特殊字符处理
func specialCharChange(url string) string {
	return strings.ReplaceAll(url, "/", "%2F")
}

func OrmTracingIsOpen() bool {
	return config.GetValueBoolDefault("base.tracing.enable", true) && config.GetValueBoolDefault("base.tracing.orm.enable", false)
}
