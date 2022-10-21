package test

import (
	"context"
	redis2 "github.com/isyscore/isc-gobase/extend/redis"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	// 客户端获取
	rdb, err := redis2.NewClient()
	if err != nil {
		logger.Warn("连接redis错误 %v", err)
		return
	}

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
