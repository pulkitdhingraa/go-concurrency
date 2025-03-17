package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// process txn with random delay
func processTransaction(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	processingTime := time.Duration(rand.Intn(7)+1) * time.Second
	done := make(chan bool)

	go func() {
		fmt.Printf("Processing transaction %d (ETA: %v)\n", id, processingTime)
		// Simulate working time - Synchronous
		time.Sleep(processingTime)
		done <- true
	}()

	// Timeout if transaction takes more than 5 seconds
	select {
	case <-done:
		fmt.Printf("Transaction %d completed succesfully\n", id)
	case <-time.After(time.Second * 5):	// timeout
		fmt.Printf("Transaction %d has timed out\n", id)
	}
}

func main() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Second * 2)	// ticker
	systemShutdown := time.NewTimer(time.Second * 15)	// timer
	defer ticker.Stop()

	fmt.Println("Payment system started")
	transactionId := 1

	Loop: 
		for {
			select {
			case <- ticker.C:
				wg.Add(1)
				go processTransaction(transactionId, &wg)
				transactionId++
			case <- systemShutdown.C:
				fmt.Println("System shutdown initiated")
				break Loop
			}
		}

	wg.Wait()	// wait for all txn to finish
	fmt.Println("All transactions processed. System shut down")
}