package main

import (
	"sync"
	"fmt"
	"context"
	"time"
)

func process(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-time.After(time.Second * 1):	// Simulate processing
		fmt.Println("Task Completed")
	case <-ctx.Done():	// Cancel if timeout occurs
		fmt.Println("Task timed out:",ctx.Err())
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	var wg sync.WaitGroup
	defer cancel()
	wg.Add(1)
	go process(ctx, &wg)
	wg.Wait()
	fmt.Println("Main func exiting")
}