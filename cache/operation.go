package cache

import "time"

//Set Add an item to the cache,replacing any existing item.
//note key is primary key
func (c *cache) Set(key string, value V) {
	de := c.defaultExpiration
	var e int64
	if de != -1 {
		e = time.Now().Add(de).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Item{
		Data: value,
		Ttl:  e,
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
	}
}

func (c *cache) Cap() {

}
