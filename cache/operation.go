package cache

import (
	"time"
)

type optList interface {
	AddItem(key string, value ...any) error
	SetItem(key string, index int, value any) error
	GetItem(key string) []any
	GetItemByIndex(key string, index int) any
	RemoveItem(key string, index int) error
}

func (c *Cache) getUnixNano() int64 {
	de := c.defaultExpiration
	var e int64
	if de != -1 {
		e = time.Now().Add(de).UnixNano()
	}
	return e
}

//Set Add an item to the cache,replacing any existing item.
//note key is primary key
func (c *Cache) Set(key string, value any) error {
	e := c.getUnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Item{
		Data: value,
		Ttl:  e,
	}
	return nil
}

//Get an item from the cache.Returns the item or nil, and a bool indicating
// whether the key was found
func (c *Cache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; !found {
		return nil, false
	} else {
		//check item has expired
		if item.Ttl > 0 && time.Now().UnixNano() > item.Ttl {
			return nil, false
		}
		return item.Data, true
	}
}

func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.items[key]; found {
		delete(c.items, key)
		return
	}
}

func (c *Cache) Clean() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.items {
		delete(c.items, k)
	}
}

func (c *Cache) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	ci := c.items
	return len(ci)
}
