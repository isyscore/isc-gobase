package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/time"
	"testing"
	t0 "time"
)

func TestWatcher(t *testing.T) {
	debug.Init()
	debug.AddWatcher("debug.test", func(key string, value string) {
		fmt.Println("最新的值1 key=", key, ", value=", value)
	})
	debug.AddWatcher("debug.test", func(key string, value string) {
		fmt.Println("最新的值2 key=", key, ", value=", value)
	})
	debug.StartWatch()

	t0.Sleep(1000000000000)
}

func TestPush(t *testing.T) {
	debug.Init()
	var tim = time.TimeToStringYmdHmsS(time.Now())
	fmt.Println(tim)
	debug.Update("debug.test", tim)
}
