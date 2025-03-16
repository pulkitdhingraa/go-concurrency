package main

import (
	"fmt"
	"sync"
	"time"
)

// wait for all queries before returning result to the user
// Reduces blocking issues by running tasks concurrently while ensuring synchronization.

func queryDatabase(db string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Querying", db)
	time.Sleep(time.Second) // Simulating DB query
	fmt.Println("Finished querying", db)
}

func main() {
	databases := []string{
		"UserDB",
		"OrderDB",
		"InventoryDB",
	}

	var wg sync.WaitGroup
	for _, db := range databases {
		wg.Add(1)
		go queryDatabase(db, &wg)
	}

	wg.Wait()
	fmt.Println("All databases queries completed")
}