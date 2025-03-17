package main

import (
	"sync"
	"fmt"
	"context"
	"time"
)

func queryDB(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		fmt.Println("Query Successful")
	case <-time.After(time.Second * 5):	// if query exceeds this deadline; return
		fmt.Println("Query exceeded deadline", ctx.Err())
	}
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second * 6))	// query with deadline time set
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go queryDB(ctx, &wg)
	wg.Wait()
	fmt.Println("Exiting main func")
}