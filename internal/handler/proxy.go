package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"wlc-loadbalancer/internal/balancer"
)

// ProxyHandler handles HTTP requests and proxies them to the selected backend
type ProxyHandler struct {
	LB *balancer.LoadBalancer
}

// NewProxyHandler creates a new proxy handler with the given load balancer
func NewProxyHandler(lb *balancer.LoadBalancer) *ProxyHandler {
	return &ProxyHandler{
		LB: lb,
	}
}

// ServeHTTP handles HTTP requests and forwards them to the selected backend
func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Select a backend using the load balancer
	backend := ph.LB.NextBackend()
	if backend == nil {
		http.Error(w, "No available backends", http.StatusServiceUnavailable)
		return
	}

	// Track incoming request
	backend.IncrementConnections()
	log.Printf("Forwarding request to %s (active connections: %d)", backend.URL.String(), backend.Connections)

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// Configure the director to modify the request before sending it to the backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = backend.URL.Host // Ensure the host header is set correctly
	}

	// Handle proxy errors
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		backend.DecrementConnections()
		http.Error(rw, "Service Unavailable", http.StatusServiceUnavailable)
	}

	// Modify the response after it's received from the backend
	proxy.ModifyResponse = func(resp *http.Response) error {
		backend.DecrementConnections() // Track completed request
		log.Printf("Request completed to %s (remaining connections: %d)",
			backend.URL.String(), backend.Connections)
		return nil
	}

	// Serve the request
	proxy.ServeHTTP(w, r)
}
