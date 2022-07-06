package test

import (
	"context"
	"github.com/isyscore/isc-gobase/goid/rl"
	"testing"
	"time"
)

type myCtxKey struct {
}

func doSomething(t *testing.T) {
	ctx := rl.Get()
	v := ctx.Value(myCtxKey{})
	t.Logf("ctx: %v\n", v)
}

func TestRoutineLocal(t *testing.T) {
	ctx := context.WithValue(context.Background(), myCtxKey{}, "233333")
	go func(t0 *testing.T) {
		rl.Set(ctx)
		time.Sleep(time.Second)
		go func(t1 *testing.T) {
			doSomething(t1)
		}(t0)
	}(t)
	time.Sleep(time.Second)
}
