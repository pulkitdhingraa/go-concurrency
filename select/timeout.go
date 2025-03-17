package main

import (
	"fmt"
	"time"
)

// timeout handling in api request
// ensures no go routine hangs forever 
// Implmenting ciruit breakers

func apicall(response chan string) {
	time.Sleep(time.Second * 3)
	response <- "API Response"
}

func main() {
	response := make(chan string)
	go apicall(response)

	select {
	case res := <-response:
		fmt.Println("Received: ", res)
	case <-time.After(time.Second * 2):
		fmt.Println("Request timed out")
	}
}
