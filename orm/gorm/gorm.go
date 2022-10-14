package gorm

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/orm"
	"github.com/isyscore/isc-gobase/tracing/sql/tracingGorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetGormDb() (*gorm.DB, error) {
	return GetGormDbWitConfig(&gorm.Config{})
}

func GetGormDbWitConfig(gormConfig *gorm.Config) (*gorm.DB, error) {
	datasourceConfig := orm.DatasourceConfig{}
	err := config.GetValueObject("base.database", &datasourceConfig)
	if err != nil {
		logger.Warn("读取datasource配置异常")
		return nil, err
	}

	var gormDb *gorm.DB
	gormDb, err = gorm.Open(mysql.Open(getDsn(datasourceConfig.DriverName)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if OrmIsOpen() {
		err := gormDb.Use(tracingGorm.NewDefault(
			config.GetValueString("base.application.name"),
			config.GetValueStringDefault("base.tracing.collector-endpoint", "http://isc-core-back-service:31300/api/core/back/v1/middle/spans"),
		))
		if err != nil {
			logger.Warn("接入tracing异常：%v", err.Error())
		}
	}

	d, _ := gormDb.DB()
	//base:
	//  datasource:
	//    username: xxx
	//    password: xxx
	//    host: xxx
	//    port: xxx
	//    # 目前支持: mysql、postgresql、sqlite、sqlserver
	//    driver-name: xxx
	//    path: xxx.sql
	//    # 示例：charset=utf8&parseTime=True&loc=Local 等url后面的配置，直接配置即可
	//    url-config:
	//      xxx: xxx
	//      yyy: yyy
	//    # 连接池配置
	//    connect-pool:
	//      # 最大空闲连接数
	//      max-idel-conns: 10
	//      # 最大连接数
	//      max-open-conns: 10
	//      # 连接可重用最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
	//      max-life-time: 10
	//      # 连接空闲的最大时间；默认单位秒，也可带字符（s：秒，m：分钟，h：小时，d：天，其他不识别按s来算）
	//      max-idle-time: 10


	//Set Max open conns设置数据库打开连接的最大数量
	d.SetMaxOpenConns()
	//Set Max idle conns设置空闲的最大连接数
	d.SetMaxIdleConns()
	//设置conn Max lifetime设置连接可重复使用的最大时间
	d.SetConnMaxLifetime()
	//设置conn最大空闲时间设置连接空闲的最大时间
	d.SetConnMaxIdleTime()

	return gormDb, nil
}

func getDsn(dbType string, datasourceConfig orm.DatasourceConfig) string {
	if dbType == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.Host, datasourceConfig.Port, datasourceConfig.DbName, datasourceConfig.SuffixUrl)
	} else if  {

	}

	switch dbType {
	case "mysql":
		break
	case "postgresql":
		break
	case "sqlite":
		break
	case "sqlserver":
		break
	}
}

func OrmIsOpen() bool {
	return config.GetValueBoolDefault("base.tracing.enable", true) && config.GetValueBoolDefault("base.tracing.orm.enable", true)
}
