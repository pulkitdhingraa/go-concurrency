package main

import (
	"fmt"
	"net/http"
)

// Handling multiple client requests in a web server

func handler(w http.ResponseWriter, r *http.Request) {
	go processRequest(r)	// non-blocking go routine processing for each request
	fmt.Fprintln(w, "Request is being processed")
}

func processRequest(r *http.Request) {
	// Simulate processing request
	fmt.Println("Processing Request:", r.URL.Path)
}


func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}