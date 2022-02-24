package test

import (
	"testing"

	"github.com/isyscore/isc-gobase/server"
)

func TestServer(t *testing.T) {
	server.InitServer()

	/*
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

	*/

	server.StartServer()
}
