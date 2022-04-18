package cache

import (
	"errors"
	"reflect"
)

//SetHash Add an item to the cache,replacing any existing item.
//note key and subKey is primary key
func (c *cache) SetHash(key, subKey string, value any) error {
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
	return nil
}

func (c *cache) RemoveHash(key, subKey string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; found {
		data := item.Data
		if reflect.TypeOf(data).Kind() != reflect.Map {
			return errors.New("key的值不是hash")
		}
		subValue := data.(map[string]any)
		delete(subValue, subKey)
		if len(subValue) == 0 {
			delete(c.items, key)
		}
		return nil
	}
	return nil
}

//GetHash get a hash value from the cache.Returns the hashes or nil, and a bool indicating
// whether the key was found
func (c *cache) GetHash(key, subKey string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; !found {
		return nil, found
	} else {
		data := item.Data
		if reflect.TypeOf(data).Kind() != reflect.Map {
			return nil, false
		}
		if value, found := item.Data.(map[string]any)[subKey]; !found {
			return nil, found
		} else {
			return value, true
		}
	}
}
