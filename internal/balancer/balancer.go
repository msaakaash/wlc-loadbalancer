package balancer

import (
	"log"
	"net/url"
	"sync"
)

// Backend represents a backend server with its associated metadata
type Backend struct {
	URL         *url.URL
	Weight      int
	Connections int // Current active connections
	mutex       sync.Mutex
}

// LoadBalancer manages a collection of backends
type LoadBalancer struct {
	Backends []*Backend
	mutex    sync.Mutex
}

// NewLoadBalancer creates a new load balancer with the given backend URLs and weights
func NewLoadBalancer(backendMap map[string]int) (*LoadBalancer, error) {
	urls := make([]string, 0, len(backendMap))
	weights := make([]int, 0, len(backendMap))

	for u, w := range backendMap {
		urls = append(urls, u)
		weights = append(weights, w)
	}

	backends := make([]*Backend, len(urls))
	for i, rawURL := range urls {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		backends[i] = &Backend{
			URL:         parsedURL,
			Weight:      weights[i],
			Connections: 0,
		}
	}
	return &LoadBalancer{Backends: backends}, nil
}

// NextBackend selects the next backend based on weighted least connections algorithm
func (lb *LoadBalancer) NextBackend() *Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	var selected *Backend
	minScore := float64(1<<63 - 1) // large number

	log.Printf("Selecting backend from:")
	for i, b := range lb.Backends {
		b.mutex.Lock()
		// If weight is 0, avoid division by zero by setting to 1
		weight := b.Weight
		if weight <= 0 {
			weight = 1
		}
		score := float64(b.Connections) / float64(weight)
		log.Printf("  Backend %d (%s): connections=%d, weight=%d, score=%.3f",
			i, b.URL.String(), b.Connections, b.Weight, score)
		b.mutex.Unlock()

		if score < minScore || (score == minScore && selected == nil) {
			minScore = score
			selected = b
		}
	}

	log.Printf("Selected backend: %s (score: %.3f)", selected.URL.String(), minScore)
	return selected
}

// IncrementConnections increases the connection count for a backend
func (b *Backend) IncrementConnections() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Connections++
}

// DecrementConnections decreases the connection count for a backend
func (b *Backend) DecrementConnections() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.Connections > 0 {
		b.Connections--
	}
}
