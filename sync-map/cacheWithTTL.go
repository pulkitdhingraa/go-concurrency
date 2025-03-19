package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cache struct {
	data sync.Map
	ttl  time.Duration
}

type cacheItem struct {
	value	interface{}
	expiryTime	time.Time
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{ttl: ttl}
	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.data.Store(key, cacheItem{
		value: value,
		expiryTime: time.Now().Add(c.ttl),
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}
	ci := item.(cacheItem)
	if time.Now().After(ci.expiryTime) {
		c.data.Delete(key)
		return nil, false
	}
	return ci.value, true
}

func (c *Cache) cleanup() {
	for {
		time.Sleep(c.ttl)
		c.data.Range(func(key, value interface{}) bool {
			ci := value.(cacheItem)
			if time.Now().After(ci.expiryTime) {
				c.data.Delete(key)
			}
			return true
		})
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	cache := NewCache(2 * time.Second)
	var wg sync.WaitGroup

	// Concurrent producers
	for i:=0;i<5;i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("Key%d",id)
			cache.Set(key, fmt.Sprintf("Value%d",id))
			fmt.Printf("Producer set %s\n",key)
		}(i)
	}

	wg.Wait()

	// Concurrent Consumers
	for i:=0;i<5;i++ {
		go func(id int) {
			key := fmt.Sprintf("Key%d",id)
			time.Sleep(time.Millisecond * 500 * time.Duration(rand.Intn(4)+1))
			val, ok := cache.Get(key)
			if ok {
				fmt.Printf("Consumer got %s = %v\n", key, val)
			} else {
				fmt.Printf("Consumer found %s expired or not present\n", key)
			}
		}(i)
	}

	// Wait to observe cleanup
	time.Sleep(time.Second * 4)
}