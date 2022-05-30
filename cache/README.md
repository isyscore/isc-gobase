# isc-gobase/cache
cache is an in-memory key:value store/cache similar to memcached that is
suitable for applications running on a single machine.

Any object can be stored, for a given duration or forever, and the cache can be
safely used by multiple goroutines.

# Installation
`go get go get github.com/isyscore/isc-gobase`

# Usage
```go
import (
	"github.com/isyscore/isc-gobase/cache"
)

func main() {
    // Create a cache with a default expiration time of forever
	c := cache.New()
	// Create a cache with an expiration time 
	c1 :=  NewWithExpiration(1 * time.Second)
	c1.Set("foo1","Vector")
	// Create a cache with an expiration time , and which
    // purges expired items every cleanup time
	c2 := NewWithExpirationAndCleanupInterval(1 * time.Second,1 * time.Second)
	c2.Set("foo2","Vector")
	
    // Set the value of the key "foo" to "bar"
	c.Set("foo", "bar")
    // Set the value of the key "baz" to 42
	c.Set("baz",42)
	// Set the value of the key "struct" to struct
    v1 := struct {
        Name string
        Age int
        }{
            "Vector.Ku",
            29,
        }
	c.Set("struct",v1)
	//Add a item to a lit
	c.AddItem("testList", "item1")
	c.AddItem("testList","item2")
	//Also,the cache allow you set an item with index,it will replace the item to new one.
	c.SetItem("testList","item3",2)
	//Get the string associated with the key "foo" from the cache
	if foo, found := c.Get("foo");found {
	    fmt.Println(foo)	
    }
	//Get the list associated with the key "testList" from the cache
	ret := c.GetItem("testList")
	for _,item := range ret {
		fmt.Println(item)
    }
}
```
