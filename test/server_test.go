package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"testing"

	"github.com/isyscore/isc-gobase/server"
)

func TestServer(t *testing.T) {
	server.InitServer()
	server.RegisterCustomHealthCheck("/api/sample",
		func() string {
			return "OK"
		},
		func() string {
			return "OK"
		},
		func() string {
			return "OK"
		},
	)

	fmt.Println(config.GetValueString("base.application.name"))
	server.StartServer()
}
