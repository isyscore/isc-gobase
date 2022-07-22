package debug

import (
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/server/rsp"
	"strings"
)

func Help(c *gin.Context) {
	port := config.GetValueIntDefault("base.server.port", 8080)
	cmdMap := map[string]any{}
	cmdMap["-============================================================================================================================================================================================"] = ""
	cmdMap["1.-【帮助】"] = "---------------------: " + "curl http://localhost:" + pre(port) + "/debug/help"
	cmdMap["2.-"] = "========================【日志】==========================================================================================================================================================="
	cmdMap["2.1 动态修改日志"] = "-----------------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"base.logger.level\", \"value\":\"debug\"}'"
	cmdMap["3.-"] = "===================【接口出入参】==========================================================================================================================================================="
	cmdMap["3.1 指定url打印请求"] = "--------------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"base.server.request.print.include-uri[0]\", \"value\":\"/api/xx/xxx\"}'"
	cmdMap["3.2 指定url不打印请求"] = "------------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"base.server.request.print.exclude-uri[0]\", \"value\":\"/api/xx/xxx\"}'"
	cmdMap["3.3 指定url打印请求和响应"] = "---------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"base.server.response.print.include-uri[0]\", \"value\":\"/api/xx/xxx\"}'"
	cmdMap["3.4 指定url不打印请求和响应"] = "-------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"base.server.response.print.exclude-uri[0]\", \"value\":\"/api/xx/xxx\"}'"
	cmdMap["4.-"] = "===================【bean管理】============================================================================================================================================================"
	cmdMap["4.1 获取注册的所有bean"] = "-----------: " + "curl http://localhost:" + pre(port) + "/bean/name/all"
	cmdMap["4.2 查询注册的bean"] = "--------------: " + "curl http://localhost:" + pre(port) + "/bean/name/list/{name}"
	cmdMap["4.3 查询bean的属性值"] = "-------------: " + "curl -X POST http://localhost:" + pre(port) + "/bean/field/get' -d '{\"bean\": \"xx\", \"field\":\"xxx\"}'"
	cmdMap["4.4 修改bean的属性值"] = "-------------: " + "curl -X POST http://localhost:" + pre(port) + "/bean/field/set' -d '{\"bean\": \"xx\", \"field\": \"xxx\", \"value\": \"xxx\"}'"
	cmdMap["4.5 调用bean的函数"] = "--------------: " + "curl -X POST http://localhost:" + pre(port) + "/bean/fun/call' -d '{\"bean\": \"xx\", \"fun\": \"xxx\", \"parameter\": {\"p1\":\"xx\", \"p2\": \"xxx\"}}'"
	cmdMap["5.-"] = "=====================【pprof】============================================================================================================================================================"
	cmdMap["5.1 动态启用pprof"] = "---------------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"base.server.gin.pprof.enable\", \"value\":\"true\"}'"
	cmdMap["6.-"] = "===================【配置处理】============================================================================================================================================================="
	cmdMap["6.1 服务所有配置"] = "----------------: " + "curl http://localhost:" + pre(port) + "/config/values"
	cmdMap["6.2 服务所有配置(yaml结构)"] = "-------: " + "curl http://localhost:" + pre(port) + "/config/values/yaml"
	cmdMap["6.3 服务某个配置"] = "----------------: " + "curl http://localhost:" + pre(port) + "/config/value/{key}"
	cmdMap["6.4 修改服务的配置"] = "--------------: " + "curl -X PUT http://localhost:" + pre(port) + "/config/update -d '{\"key\":\"xxx\", \"value\":\"yyy\"}'"
	cmdMap["==============================================================================================================================================================================================="] = ""

	rsp.Success(c, cmdMap)
}

func pre(port int) string {
	return isc.ToString(port) + apiPreAndModule()
}

func apiPreAndModule() string {
	apiPrefix := config.GetValueStringDefault("base.api.prefix", "/api")
	apiPrefix = strings.TrimSuffix(apiPrefix, "/")
	module := strings.Trim(config.ApiModule, "/")
	if module == "" {
		return apiPrefix
	}

	return apiPrefix + "/" + module
}
