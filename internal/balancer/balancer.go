package balancer

import (
	"net/url"
	"sync"
)

type Backend struct {
	URL         *url.URL
	Weight      int
	Connections int
}

type LoadBalancer struct {
	Backends []*Backend
	mutex    sync.Mutex
}

func NewLoadBalancer(urls []string, weights []int) *LoadBalancer {
	backends := make([]*Backend, len(urls))
	for i, rawURL := range urls {
		parsedURL, _ := url.Parse(rawURL)
		backends[i] = &Backend{
			URL:         parsedURL,
			Weight:      weights[i],
			Connections: 0,
		}
	}
	return &LoadBalancer{Backends: backends}
}

func (lb *LoadBalancer) NextBackend() *Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	var selected *Backend
	minScore := 1e9 // large float

	for _, b := range lb.Backends {
		score := float64(b.Connections) / float64(b.Weight)
		if score < minScore {
			minScore = score
			selected = b
		}
	}

	return selected
}

func (lb *LoadBalancer) IncrementConnections(backend *Backend) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	backend.Connections++
}

func (lb *LoadBalancer) DecrementConnections(backend *Backend) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	if backend.Connections > 0 {
		backend.Connections--
	}
}
