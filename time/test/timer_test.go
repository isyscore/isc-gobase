package test

import (
	"testing"

	t0 "time"

	"github.com/isyscore/isc-gobase/time"
)

func TestTimer(t *testing.T) {

	globalCount := 0
	timer := time.NewTimerWithFire(1.5, func(tm *time.Timer) {
		// 这里是运行在协程内的
		globalCount++
		t.Logf("globalCount: %d", globalCount)
	})
	timer.Start()
	t.Logf("timer.stopped = %v\n", timer.IsStopped)

	// 等待timer执行5次
	for globalCount != 5 {
		t0.Sleep(t0.Second)
	}

	timer.Stop()
	t.Logf("timer.stopped = %v\n", timer.IsStopped)

	// 证明timer已停
	for globalCount != 10 {
		t0.Sleep(t0.Second)
		globalCount++
	}
}

func TestTimerParam(t *testing.T) {
	globalCount := 0
	timer := time.NewTimerWithFire(1.5, func(tm *time.Timer) {
		// 这里是运行在协程内的
		globalCount++
		t.Logf("globalCount: %d", globalCount)
	})
	timer.Start()
	t.Logf("timer.stopped = %v\n", timer.IsStopped)
	// 等待timer执行3次
	for globalCount != 3 {
		t0.Sleep(t0.Second)
	}
	// 修改参数
	timer.SetInterval(3)
	timer.SetOnTimer(func(tm *time.Timer) {
		globalCount++
		t.Logf("globalCount2: %d", globalCount)
	})
	// 等待timer执行3次
	for globalCount != 6 {
		t0.Sleep(t0.Second)
	}
	timer.Stop()
}
