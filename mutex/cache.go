package main

import (
	"fmt"
	"sync"
)

// Shared resource to be r/w by multiple goroutines

type Cache struct {
	data map[string]string
	mu   sync.RWMutex
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, exists := c.data[key]
	return val, exists
}

func main() {
	cache := Cache{
		data: make(map[string]string),
	}

	var wg sync.WaitGroup

	// 5 go routines writing to cache concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			cache.Set(key, fmt.Sprintf("value-%d", id))
			fmt.Printf("Routine %d set key to: %s\n", id, key)
		}(i)
	}

	// 5 go routines reading from cache concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			val, exists := cache.Get(key)
			if exists {
				fmt.Printf("Reader %d got key: %s with value: %s\n", id, key, val)
			} else {
				fmt.Printf("Reader %d could not find key: %s\n", id, key)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("All go routines completed")
}
