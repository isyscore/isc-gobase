package orm

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/sirupsen/logrus"
	"strings"
)

func getDbDsnWithName(datasourceName string) (string, error) {
	datasourceConfig := config.DatasourceConfig{}
	targetDatasourceName := "base.datasource"
	if datasourceName != "" {
		targetDatasourceName = "base.datasource." + datasourceName
	}
	err := config.GetValueObject(targetDatasourceName, &datasourceConfig)
	if err != nil {
		logger.Warn("读取读取配置【datasource】异常")
		return "", err
	}

	return getDbDsn(datasourceConfig.DriverName, datasourceConfig), nil
}

func getDbDsn(dbType string, datasourceConfig config.DatasourceConfig) string {
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
		return dsn
	case "postgresql":
		// 格式：host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
		// host=127.0.0.1 port=54321 user=isyscore password=Isysc0re dbname=isyscore sslmode=disable search_path=isc_permission
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
			datasourceConfig.Host, datasourceConfig.Port, datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.Username, datasourceConfig.DbName)
		//if len(sqlConfigMap) != 0 {
		//	var kvList []string
		//	for key, value := range sqlConfigMap {
		//		kvList = append(kvList, fmt.Sprintf("%s=%s", key, value))
		//	}
		//	dsn += fmt.Sprintf(" %s", strings.Join(kvList, " "))
		//}
		return dsn
	case "sqlite":
		// 格式： gorm.db
		return datasourceConfig.SqlitePath
	case "sqlserver":
		// 格式：sqlserver://user:password@localhost:9930?database=gorm
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.Host, datasourceConfig.Port, datasourceConfig.DbName)
		return dsn
	}
	return ""
}

func ConfigChangeListenerOfOrm(event listener.BaseEvent) {
	ev := event.(listener.ConfigChangeEvent)
	if ev.Key == "base.orm.show-sql" {
		if isc.ToBool(ev.Value) {
			logger.Group("orm").SetLevel(logrus.DebugLevel)
		} else {
			logger.Group("orm").SetLevel(logrus.InfoLevel)
		}
	}
}
