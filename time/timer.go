package time

import (
	t0 "time"
)

type OnTimerFunc func(t *Timer)

type Timer struct {
	Interval      int64
	IsStopped     bool
	OnTimer       OnTimerFunc // 这个函数将在[协程]内异步调用
	OnBeforeTimer OnTimerFunc // 这个函数在[主线程]内同步调用
	OnAfterTimer  OnTimerFunc // 这个函数在[主线程]内同步调用
	ticker        *t0.Ticker
	shiftStatus   bool
}

func NewTimer() *Timer {
	intv := int64(float64(t0.Second) * 1.0)
	return &Timer{
		Interval:    intv,
		IsStopped:   true,
		ticker:      nil,
		shiftStatus: false,
	}
}

func NewTimerWithInterval(seconds float64) *Timer {
	intv := int64(float64(t0.Second) * seconds)
	return &Timer{
		Interval:    intv,
		IsStopped:   true,
		ticker:      nil,
		shiftStatus: false,
	}
}

func NewTimerWithFire(seconds float64, onTimer OnTimerFunc) *Timer {
	intv := int64(float64(t0.Second) * seconds)
	return &Timer{
		Interval:    intv,
		IsStopped:   true,
		OnTimer:     onTimer,
		ticker:      nil,
		shiftStatus: false,
	}
}

func (t *Timer) Start() {
	t.IsStopped = false
	if !t.shiftStatus {
		if t.OnBeforeTimer != nil {
			t.OnBeforeTimer(t)
		}
	}

	t.ticker = t0.NewTicker(t0.Duration(t.Interval))
	go func(tk *t0.Ticker) {
		for !t.IsStopped {
			_, ok := <-tk.C
			if !ok {
				break
			}
			if t.OnTimer != nil {
				t.OnTimer(t)
			}
		}
	}(t.ticker)
}

func (t *Timer) Stop() {
	t.IsStopped = true
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
	if !t.shiftStatus {
		if t.OnAfterTimer != nil {
			t.OnAfterTimer(t)
		}
	}
}

func (t *Timer) SetInterval(seconds float64) {
	shouldStart := false
	t.shiftStatus = true
	if !t.IsStopped {
		shouldStart = true
		t.Stop()
	}
	intv := int64(float64(t0.Second) * seconds)
	t.Interval = intv
	if shouldStart {
		t.Start()
	}
	t.shiftStatus = false
}

func (t *Timer) SetOnTimer(event OnTimerFunc) {
	shouldStart := false
	t.shiftStatus = true
	if !t.IsStopped {
		shouldStart = true
		t.Stop()
	}
	t.OnTimer = event
	if shouldStart {
		t.Start()
	}
	t.shiftStatus = false
}

func (t *Timer) SetOnBefore(event OnTimerFunc) {
	t.OnBeforeTimer = event
}

func (t *Timer) SetOnAfter(event OnTimerFunc) {
	t.OnAfterTimer = event
}
