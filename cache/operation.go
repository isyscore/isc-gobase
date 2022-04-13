package cache

import (
	"time"
)

func (c *cache) getUnixNano() int64 {
	de := c.defaultExpiration
	var e int64
	if de != -1 {
		e = time.Now().Add(de).UnixNano()
	}
	return e
}

func (c *cache) getLock() {
	for !c.mu.TryLock() {
		time.Sleep(10 * time.Millisecond)
	}
}

//Set Add an item to the cache,replacing any existing item.
//note key is primary key
func (c *cache) Set(key string, value any) {
	e := c.getUnixNano()
	c.getLock()
	defer c.mu.Unlock()
	c.items[key] = Item{
		Data: value,
		Ttl:  e,
	}
}

//Get an item from the cache.Returns the item or nil, and a bool indicating
// whether the key was found
func (c *cache) Get(key string) (any, bool) {
	c.getLock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; !found {
		return nil, found
	} else {
		//check item has expired
		if item.Ttl > 0 && time.Now().UnixNano() > item.Ttl {
			return nil, found
		}
		return item.Data, true
	}
}

func (c *cache) Remove(key string) {
	c.getLock()
	defer c.mu.Unlock()
	if _, found := c.items[key]; found {
		delete(c.items, key)
		return
	}
}

func (c *cache) Cap() int {
	c.getLock()
	defer c.mu.Unlock()
	ci := c.items
	return len(ci)
}
