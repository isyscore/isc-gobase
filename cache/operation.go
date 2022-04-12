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

//Set Add an item to the cache,replacing any existing item.
//note key is primary key
func (c *cache) Set(key string, value V) {
	e := c.getUnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Item{
		Data: value,
		Ttl:  e,
	}
}

type HashStruct struct {
	Key  string
	Data V
}

//SetHash Add an item to the cache,replacing any existing item.
//note key and subKey is primary key
func (c *cache) SetHash(key, subKey string, value V) {
	e := c.getUnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	if vv, b := c.items[key]; b {
		vv.Data.(map[string]any)[subKey] = value
	} else {
		subKeyValue := make(map[string]any)
		subKeyValue[subKey] = value
		v := Item{
			Data: subKeyValue,
			Ttl:  e,
		}
		c.items[key] = v
	}
}

//GetHash get a hash value from the cache.Returns the hashes or nil, and a bool indicating
// whether the key was found
func (c *cache) GetHash(key, subKey string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; !found {
		return nil, found
	} else {
		if value, found := item.Data.(map[string]any)[subKey]; !found {
			return nil, found
		} else {
			return value, true
		}
	}
}

//Get an item from the cache.Returns the item or nil, and a bool indicating
// whether the key was found
func (c *cache) Get(key string) (V, bool) {
	c.mu.Lock()
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
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, found := c.items[key]; found {
		delete(c.items, key)
		return
	}
}

func (c *cache) Cap() int {
	ci := c.items
	return len(ci)
}
