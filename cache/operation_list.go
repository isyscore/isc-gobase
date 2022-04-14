package cache

import (
	"errors"
	"reflect"
)

func (c *cache) AddItem(key string, value ...any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, found := c.items[key]; found {
		data := item.Data
		if reflect.TypeOf(data).Kind() != reflect.Slice {
			return errors.New("key 对应的数据类型不是 slice")
		}
		item.Data = append(data.([]any), value)
		c.items[key] = item
	} else {
		e := c.getUnixNano()
		data := Item{
			Data: []any{value},
			Ttl:  e,
		}
		c.items[key] = data
	}
	return nil
}

//SetItem set or replace a value of items by index
func (c *cache) SetItem(key string, idx int, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; found {
		data := item.Data
		if reflect.TypeOf(data).Kind() != reflect.Slice {
			return errors.New("key 对应的数据类型不是 slice")
		}
		items := data.([]any)
		if len(items) <= idx {
			return errors.New("数组下标越界")
		}
		item.Data = append(item.Data.([]any), value)
	}
	return errors.New("key不存在")
}

//GetItem return an array of points or nil
func (c *cache) GetItem(key string) []any {
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, found := c.items[key]; !found {
		return nil
	} else {
		data := item.Data
		if reflect.TypeOf(data).Kind() != reflect.Slice {
			return nil
		}
		return data.([]any)
	}
}

//GetItemByIndex return a value of Type is T or nil
func (c *cache) GetItemByIndex(key string, idx int) any {
	if idx < 0 {
		return nil
	}
	if items := c.GetItem(key); items != nil {
		if len(items) <= idx {
			return nil
		}
		return items[idx]
	}
	return nil
}

//RemoveItem an item from the cache. Does nothing if the key is not in the cache.
func (c *cache) RemoveItem(key string, idx int) error {
	if idx < 0 {
		c.Remove(key)
		return nil
	}
	item := c.GetItem(key)
	if len(item) <= idx {
		return nil
	}
	newItem := item[:idx]
	newItem = append(newItem, item[idx:]...)
	c.Set(key, newItem)
	return nil
}
