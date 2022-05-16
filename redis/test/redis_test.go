package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/isyscore/isc-gobase/redis"
)

func TestConnect(t *testing.T) {
	rdb, _ := redis.GetClient()

	ctx := context.Background()
	rdb.Set(ctx, "k1", "vv", time.Hour)
	rlt := rdb.Get(ctx, "k1")
	fmt.Println(rlt.Result())
}
