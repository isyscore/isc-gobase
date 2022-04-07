package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/isyscore/isc-gobase/config"

	"github.com/gin-gonic/gin"
	h2 "github.com/isyscore/isc-gobase/http"
	t2 "github.com/isyscore/isc-gobase/time"
)

var procId = os.Getpid()
var startTime = time.Now().Format(t2.FmtYMdHms)

const defaultVersion = "unknown"

var Version string = defaultVersion

func healthSystemStatus(c *gin.Context) {
	c.Data(http.StatusOK, h2.ContentTypeJson, []byte(fmt.Sprintf(`{"status":"ok","running":true,"pid":%d,"startupAt":"%s","version":"%s"}`, procId, startTime, getVersion())))
}

func healthSystemInit(c *gin.Context) {
	c.Data(http.StatusOK, h2.ContentTypeText, []byte(`{"status":"ok"}`))
}

func healthSystemDestroy(c *gin.Context) {
	c.Data(http.StatusOK, h2.ContentTypeText, []byte(`{"status":"ok"}`))
}

func getVersion() string {
	if Version != defaultVersion {
		return Version
	}
	Version := config.GetValueStringDefault("base.server.version", defaultVersion)
	return Version
}
