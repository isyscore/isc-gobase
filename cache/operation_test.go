package cache

import (
	"fmt"
	"strconv"
	"testing"
	"time"
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

func Test_cache_Get1(t *testing.T) {
	c := NewWithExpiration(1 * time.Second)
	c.Set("testV", "库陈胜")
	time.Sleep(1 * time.Second)
	get, b := c.Get("testV")
	if b {
		t.Error("expire fail get is ", get)
	}
	c.SetHash("Key", "subKey", "库陈胜")
	c.SetHash("Key", "6不liuKey@", "贼溜")

	if v1, b := c.Get("Key"); !b {
		t.Error("未获取到hash值")
	} else {
		t.Logf("获取到hash值 %v", v1)
	}
	time.Sleep(1 * time.Second)

	if v2, b := c.GetHash("Key", "subKey"); !b {
		t.Error("未获取到hash值 v2")
	} else {
		t.Logf("获取到hash 及子值 %v", v2)
	}
	time.Sleep(2 * time.Second)

	if _, b = c.Get("Key"); b {
		t.Error("Key 未过期")
	}

	if _, b = c.GetHash("Key", "subKey"); b {
		t.Error("Key - subKey 未过期")
	}
}

//性能测试
//fixme 并发问题有待处理
func Test_cache_Get2(t *testing.T) {
	c := NewWithExpiration(5 * time.Second)
	length := 300000
	ch := make(chan int8, length)
	start := time.Now()
	println("开始执行", start.UnixNano())
	for i := 0; i < length/3; i++ {
		key := fmt.Sprintf("%s%d", "Key", i)
		go func(ii int, k string) {
			c.Set(k, "库陈胜"+k)
			ch <- int8(1)
			c.SetHash(k+"hash", strconv.Itoa(ii), "性能测试"+k)
			ch <- int8(1)
			c.AddItem(k, ii)
			ch <- int8(1)
		}(i, key)
	}

	for i := 0; i < length; i++ {
		<-ch
	}
	close(ch)

	//for c.Cap() != 20000 {
	//	time.Sleep(100 * time.Millisecond)
	//	t.Logf("CAP %d", c.Cap())
	//}

	t.Logf("PUT结束执行,耗时 %d ms, key总数: %d", time.Now().UnixMilli()-start.UnixMilli(), c.Cap())

	ch1 := make(chan int8, length)
	for i := 0; i < length/3; i++ {
		key := fmt.Sprintf("%s%d", "Key", i)
		subKey := key + "hash"
		go func(k, s string) {
			_, _ = c.Get(k)
			ch1 <- int8(1)
			_, _ = c.GetHash(key+"hash", s)
			ch1 <- int8(1)
			ret := c.GetItem(k)
			fmt.Println("k=", k, "value=", ret)
		}(key, subKey)
	}
	for i := 0; i < length; i++ {
		<-ch1
	}
	close(ch1)

	t.Logf("当前key的数量 = %d", c.Cap())
	times := 1
	for c.Cap() > 0 {
		time.Sleep(time.Second)
		t.Log("沉睡", times, "秒后，剩余多少Key?", c.Cap())
		times++
	}
}

func TestCache_AddItem(t *testing.T) {
	c := NewWithExpiration(5 * time.Second)
	start := time.Now().UnixMilli()
	ch := make(chan int8, 100)
	for i := 0; i < 100; i++ {
		go func(ii int) {
			if err := c.AddItem("test", ii, 666); err != nil {
				t.Errorf("%v", err)
			}
			ch <- int8(1)
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-ch
	}
	t.Log("执行结束，耗时", time.Now().UnixMilli()-start, "ms", "缓存大小", c.Cap())
	ret := c.GetItem("test")
	t.Logf("size = %d,%v", len(ret),
		ret)
}
