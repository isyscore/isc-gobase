package test

import (
	"context"
	redis2 "github.com/isyscore/isc-gobase/extend/redis"
	"github.com/magiconair/properties/assert"
	goredis "github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var rdb goredis.UniversalClient

func init() {
	// 客户端获取
	_rdb, err := redis2.NewClient()
	if err != nil {
		return
	}
	rdb = _rdb
}

func TestRedis(t *testing.T) {
	// 添加和读取
	key := "test_key"
	value := "test_value"

	ctx := context.Background()
	rdb.Set(ctx, key, value, time.Hour)
	rlt := rdb.Get(ctx, key)

	// 判断
	actValue, _ := rlt.Result()
	assert.Equal(t, actValue, value)
}
