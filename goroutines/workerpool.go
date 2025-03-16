package main

import (
	"fmt"
	"time"
)

// Limit concurrent execution with bounded no. of go routines
// Rate Limiting, Task Queues and Background Jobs
// Multiple jobs to be processed in parallel

func worker(id int, jobs chan int, results chan int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(time.Second)
		results <- job*2
	}
}

func main() {
	jobs := make(chan int, 5)	// jobs channel having capacity of 5
	results := make(chan int, 10)

	// Create 3 worker go routines to handle 10 jobs
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 10 jobs
	for j:=1;j<=10;j++ {
		jobs <- j
	}

	close(jobs)

	for a := 1; a <= 10; a++ {
		fmt.Println("Results: ", <-results)
	}
}