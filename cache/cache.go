package cache

import (
	"runtime"
	"sync"
	"time"
)

type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, any)
	j                 *janitor
}

type Item struct {
	//data
	Data any
	//ttl time.UnixNano.Expiration of Item,if it is -1,it will be not Expired
	Ttl int64
}

func (item *Item) Expired() bool {
	if item.Ttl == -1 {
		//用户
		return false
	}
	return time.Now().UnixNano() > item.Ttl
}

func New() *Cache {
	return NewWithExpiration(0)
}

func NewWithExpiration(expiration time.Duration) *Cache {
	return NewWithExpirationAndCleanupInterval(expiration, 0)
}

func NewWithCleanupInterval(cleanupInterval time.Duration) *Cache {
	return NewWithExpirationAndCleanupInterval(0, cleanupInterval)
}
func NewWithExpirationAndCleanupInterval(defaultExpiration, cleanupInterval time.Duration) *Cache {
	if defaultExpiration == 0 {
		defaultExpiration = -1
	}
	ch := make(chan bool)
	c := &Cache{
		defaultExpiration: defaultExpiration,
		items:             make(map[string]Item),
		j: &janitor{
			Interval: 1 * time.Second,
			stop:     ch,
		},
	}
	//启动清理协程
	go func() {
		c.runCleanup(cleanupInterval)
		runtime.SetFinalizer(c, stopJanitor)
	}()
	return c
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (c *Cache) runCleanup(cleanupInterval time.Duration) {
	if cleanupInterval == 0 {
		cleanupInterval = 500 * time.Millisecond
	}
	ticker := time.NewTicker(cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-c.j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.j.stop <- true
}

func (c *Cache) DeleteExpired() {
	l := len(c.items)
	if l < 1 {
		return
	}
	cloneMap := map[string]Item{}
	c.mu.Lock()
	for k, v := range c.items {
		cloneMap[k] = v
	}
	c.mu.Unlock()

	ch := make(chan int8, len(c.items))
	for key, item := range cloneMap {
		go func(k string, i Item) {
			if i.Expired() {
				c.mu.Lock()
				delete(c.items, k)
				c.mu.Unlock()
			}
			ch <- int8(1)
		}(key, item)
	}
}
