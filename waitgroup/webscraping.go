package main

import (
	"fmt"
	"net/http"
	"sync"
)

// wait for all parallel requests to complete before processing results

func fetchURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching:", url)
		return
	}

	fmt.Println("Fetched:", url, "Status:", res.Status)
}

func main() {
	urls := []string {
		"https://golang.org",
		"https://google.com",
		"https://github.com",
	}

	var wg sync.WaitGroup	// value type
	// wg := &sync.WaitGroup{}	// pointer type
	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, &wg)
	}

	wg.Wait()	// wait for all requests to complete before exiting main go routine
	fmt.Println("All URLs fetched")
}