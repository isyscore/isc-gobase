package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	h2 "github.com/isyscore/isc-gobase/http"
	t2 "github.com/isyscore/isc-gobase/time"
)

var procId = os.Getpid()
var startTime = time.Now().Format(t2.FmtYMdHms)

func healthSystemStatus(c *gin.Context) {
	c.Data(http.StatusOK, h2.ContentTypeJson, []byte(fmt.Sprintf(`{"status":"ok","running":true,"pid":%d,"startupAt":"%s"}`, procId, startTime)))
}

func healthSystemInit(c *gin.Context) {
	c.Data(http.StatusOK, h2.ContentTypeText, []byte(`{"status":"ok"}`))
}

func healthSystemDestroy(c *gin.Context) {
	c.Data(http.StatusOK, h2.ContentTypeText, []byte(`{"status":"ok"}`))
}
