package orm

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	"strings"
)

type GobaseSqlHook interface {
	Before(ctx context.Context, parameters map[string]any) (context.Context, error)
	After(ctx context.Context, parameters map[string]any) (context.Context, error)
	Err(ctx context.Context, err error, parameters map[string]any) error
}

var SqlHooks []GobaseSqlHook

func init() {
	SqlHooks = []GobaseSqlHook{}
}

func AddHook(hook GobaseSqlHook) {
	SqlHooks = append(SqlHooks, hook)
}

type GobaseSqlHookProxy struct {}

func (proxy *GobaseSqlHookProxy) Before(ctx context.Context, query string, args ...interface {}) (context.Context, error) {
	for _, hook := range SqlHooks {
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
	for _, hook := range SqlHooks {
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
	for _, hook := range SqlHooks {
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
		fmt.Println(dsn)
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
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", datasourceConfig.Host, datasourceConfig.Username, datasourceConfig.Password, datasourceConfig.DbName, datasourceConfig.Port)
		if len(sqlConfigMap) != 0 {
			var kvList []string
			for key, value := range sqlConfigMap {
				kvList = append(kvList, fmt.Sprintf("%s=%s", key, value))
			}
			dsn += fmt.Sprintf(" %s", strings.Join(kvList, " "))
		}
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

