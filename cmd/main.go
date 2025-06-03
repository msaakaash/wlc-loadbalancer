package main

import (
	"log"
	"net/http"
	"time"
	"wlc-loadbalancer/internal/balancer"
	"wlc-loadbalancer/internal/handler"
	"wlc-loadbalancer/internal/server"
)

func main() {
	// Launch backends in goroutines
	go server.Start("8081")
	go server.Start("8082")
	go server.Start("8083")

	// Wait a bit to ensure backends are up
	log.Println("Waiting for backend servers to start...")
	time.Sleep(2 * time.Second)

	// Backend config: map[URL]weight
	backends := map[string]int{
		"http://localhost:8081": 3, // Higher weight = more traffic
		"http://localhost:8082": 1,
		"http://localhost:8083": 2,
	}

	// Create the load balancer
	lb, err := balancer.NewLoadBalancer(backends)
	if err != nil {
		log.Fatal(err)
	}

	// Create the proxy handler
	proxy := handler.NewProxyHandler(lb)

	// Start the proxy server
	log.Println("Load balancer is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
