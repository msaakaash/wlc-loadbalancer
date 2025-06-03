package server

import (
	"fmt"
	"log"
	"net/http"
)

// Start launches a simple HTTP server on the specified port
func Start(port string) {
	mux := http.NewServeMux()

	// Simple handler that returns the server's port
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request on server:%s %s %s", port, r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello from backend server on port %s!", port)
	})

	// Launch the server
	addr := ":" + port
	log.Printf("Starting backend server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Server on port %s failed: %v", port, err)
	}
}
