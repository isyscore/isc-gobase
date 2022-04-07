package cache

import (
	"testing"
)

type Value struct {
	Name string
	Age  int
}

func Test_cache_Get(t *testing.T) {
	c := New()
	v := Value{"库陈胜", 2022}
	c.Set("test", v)
	c.Set("test1", "啊哈")
	get, err := c.Get("test")
	t.Logf("result %v,%v", get, err)

	c.Remove("test")
	get, err = c.Get("test")
	t.Logf("result %v,%v", get, err)
}
