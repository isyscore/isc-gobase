package test

import (
	"context"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"

	"github.com/isyscore/isc-gobase/redis"
)

func TestRedis(t *testing.T) {
	// 客户端获取
	rdb, _ := redis.GetClient()

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
