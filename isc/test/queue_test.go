package test

import (
	"fmt"
	"github.com/isyscore/isc-gobase/isc"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	queue := isc.NewQueue()

	// 返回当前还有多少个
	num := queue.Offer("dsf")
	fmt.Println(num)

	// 返回值，以及返回值之后的剩余个数
	d, n := queue.Poll()
	fmt.Println(d, n)

	// 获取数据
	d, n = queue.Take(1 * time.Second)
	fmt.Println(d, n)

	d, n = queue.Take(1 * time.Second)
	fmt.Println(d, n)
}
