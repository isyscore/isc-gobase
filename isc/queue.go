package isc

import (
	"sync"
	"sync/atomic"
	t0 "time"
)

type element struct {
	next  *element
	value interface{}
}

func newElement(value interface{}) *element {
	return &element{value: value}
}

type Queue struct {
	input   *element
	output  *element
	ch      chan bool
	waitNum int32
	num     int32
	lock    *sync.RWMutex
	rLock   *sync.Mutex
	timer   *timer
}

type timer struct {
	t     *t0.Timer
	isEnd chan bool
}

func newTimer() *timer {
	return &timer{t: t0.NewTimer(t0.Second * 10), isEnd: make(chan bool)}
}

func (timer *timer) wait() bool {
	select {
	case <-timer.t.C:
		return false
	case fa := <-timer.isEnd:
		return fa
	}
	// return false
}

func (timer *timer) end() {
	fa := timer.t.Stop()
	if fa {
		timer.isEnd <- true
	}
}

func (timer *timer) reset(duration t0.Duration) {
	timer.t.Reset(duration)
}

var poolTimer = &sync.Pool{
	New: func() interface{} {
		return newTimer()
	},
}

func getTimer(duration t0.Duration) *timer {
	ti := poolTimer.Get().(*timer)
	ti.reset(duration)
	return ti
}

func freeTimer(timer *timer) {
	poolTimer.Put(timer)
}

func NewQueue() *Queue {
	return &Queue{ch: make(chan bool), waitNum: 0, num: 0, lock: new(sync.RWMutex), rLock: new(sync.Mutex)}
}

func (queue *Queue) Offer(value interface{}) (num int32) {
	ele := newElement(value)
	queue.lock.Lock()
	if queue.num == 0 {
		queue.input = ele
		queue.output = ele
	} else {
		queue.input.next = ele
		queue.input = ele
	}
	num = atomic.AddInt32(&queue.num, 1)
	if queue.waitNum > 0 {
		queue.waitNum--
		queue.lock.Unlock()
		queue.ch <- true
	} else {
		queue.lock.Unlock()
	}
	return
}

func (queue *Queue) Num() int32 {
	return queue.num
}

func (queue *Queue) Peek() (value interface{}, num int32) {
	queue.lock.RLock()
	num = queue.num
	if queue.num > 0 {
		value = queue.output.value
		queue.lock.RUnlock()
		return
	} else {
		queue.lock.RUnlock()
		return nil, 0
	}
}

func (queue *Queue) Poll() (value interface{}, num int32) {
	for {
		queue.lock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				value, num = queue.readOne()
				queue.lock.Unlock()
				return
			} else {
				queue.lock.Unlock()
				queue.rLock.Lock()
				val, n, last := queue.readGtOne()
				if last {
					queue.rLock.Unlock()
				} else {
					queue.rLock.Unlock()
					return val, n
				}
			}
		} else {
			queue.waitNum++
			queue.lock.Unlock()
			<-queue.ch
		}
	}
}

func (queue *Queue) readOne() (value interface{}, num int32) {
	var ele = queue.output
	value = ele.value
	num = atomic.AddInt32(&queue.num, -1)
	return value, num
}

func (queue *Queue) readGtOne() (value interface{}, num int32, isLast bool) {
	var ele = queue.output
	if ele.next == nil {
		return nil, 0, true
	}
	value = ele.value
	queue.output = ele.next
	ele.next = nil
	num = atomic.AddInt32(&queue.num, -1)
	return value, num, false
}

func (queue *Queue) Take(duration t0.Duration) (value interface{}, num int32) {
	for {
		queue.lock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				value, num = queue.readOne()
				queue.lock.Unlock()
				return
			} else {
				queue.lock.Unlock()
				queue.rLock.Lock()
				val, n, last := queue.readGtOne()
				if last {
					queue.rLock.Unlock()
				} else {
					queue.rLock.Unlock()
					return val, n
				}
			}
		} else {
			queue.waitNum++
			queue.lock.Unlock()
			tm := getTimer(duration)
			go func() {
				fa := tm.wait()
				if !fa {
					queue.lock.Lock()
					if queue.waitNum > 0 {
						queue.waitNum--
						queue.lock.Unlock()
						queue.ch <- false
					} else {
						queue.lock.Unlock()
					}
				}
			}()
			flag := <-queue.ch
			tm.end()
			freeTimer(tm)
			if !flag {
				return nil, 0
			}
		}
	}
}
