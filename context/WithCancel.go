package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <- ctx.Done():
			fmt.Printf("Worker %d stopped\n", id)
			return
		default:
			fmt.Printf("Worker %d is working\n", id)
			time.Sleep(time.Second * 1)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(2)
	go worker(ctx, 1, &wg)
	go worker(ctx, 2, &wg)

	time.Sleep(time.Second * 3)
	fmt.Println("Cancelling workers...")
	// cancelling the context stops all associated go routines
	cancel()
	wg.Wait()
	fmt.Println("Main function exiting")
}