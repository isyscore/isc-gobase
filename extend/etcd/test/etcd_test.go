package test

import (
	"github.com/isyscore/isc-gobase/extend/etcd"
	"testing"
)

func Test1(t *testing.T) {
	client, _ := etcd.NewEtcdClient()
	client.Put()
}
