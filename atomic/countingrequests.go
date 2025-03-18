package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var activeRequests int32 = 0

// Simulate requests increasing over time
func handleRequest(wg *sync.WaitGroup){
	defer wg.Done()
	atomic.AddInt32(&activeRequests, 2)
	time.Sleep(time.Millisecond * 500)
	atomic.AddInt32(&activeRequests, -1)
}

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	for i:=0; i<10; i++ {
		wg.Add(1)
		go handleRequest(&wg)
	}

	// Monitor active requests
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("✅ Monitor exiting due to context cancellation")
				return
			default:
				current := atomic.LoadInt32(&activeRequests)
				fmt.Printf("Active requests: %d\n",current)
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()

	wg.Wait()
	cancel()
	fmt.Println("🏁 All done, exiting.")
}