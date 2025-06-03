package balancer

import (
	"testing"
)

func TestNewLoadBalancer(t *testing.T) {
	backends := map[string]int{
		"http://backend1.com": 2,
		"http://backend2.com": 1,
	}

	lb, err := NewLoadBalancer(backends)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(lb.Backends) != 2 {
		t.Errorf("Expected 2 backends, got %d", len(lb.Backends))
	}
}

func TestNextBackend(t *testing.T) {
	backends := map[string]int{
		"http://backend1.com": 2,
		"http://backend2.com": 1,
		"http://backend3.com": 3,
	}

	lb, _ := NewLoadBalancer(backends)
	lb.Backends[0].Connections = 4
	lb.Backends[1].Connections = 3
	lb.Backends[2].Connections = 6

	selected := lb.NextBackend()
	selectedURL := selected.URL.String()

	if selectedURL != "http://backend1.com" && selectedURL != "http://backend3.com" {
		t.Errorf("Expected http://backend1.com or http://backend3.com, but got %s", selectedURL)
	}
}

func TestIncrementConnections(t *testing.T) {
	backends := map[string]int{
		"http://backend1.com": 1,
	}

	lb, _ := NewLoadBalancer(backends)
	backend := lb.Backends[0]
	lb.IncrementConnections(backend)

	if backend.Connections != 1 {
		t.Errorf("Expected 1 connection, got %d", backend.Connections)
	}
}

func TestDecrementConnections(t *testing.T) {
	backends := map[string]int{
		"http://backend1.com": 1,
	}

	lb, _ := NewLoadBalancer(backends)
	backend := lb.Backends[0]
	backend.Connections = 2
	lb.DecrementConnections(backend)

	if backend.Connections != 1 {
		t.Errorf("Expected 1 connection, got %d", backend.Connections)
	}
}

func TestNextBackendWithEqualConnections(t *testing.T) {
	backends := map[string]int{
		"http://backend1.com": 2,
		"http://backend2.com": 1,
		"http://backend3.com": 3,
	}

	lb, _ := NewLoadBalancer(backends)
	for _, backend := range lb.Backends {
		backend.Connections = 3
	}

	selected := lb.NextBackend()

	if selected.URL.String() != "http://backend3.com" {
		t.Errorf("Expected http://backend3.com (highest weight), got %s", selected.URL.String())
	}
}
