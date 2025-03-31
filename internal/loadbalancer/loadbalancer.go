package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gentrit-hoxha/golang-load-balancer/internal/jwt"
)

// LoadBalancer structure
type LoadBalancer struct {
	servers      []*url.URL
	currentIdx   int
	mu           sync.Mutex
	jwtValidator *jwt.Validator
}

// Create new load balancer
func NewLoadBalancer(serverURLs []string, jwtValidator *jwt.Validator) *LoadBalancer {
	parsedServers := make([]*url.URL, len(serverURLs))
	for i, serverURL := range serverURLs {
		parsed, err := url.Parse(serverURL)
		if err != nil {
			log.Fatalf("Invalid server URL: %v", err)
		}
		parsedServers[i] = parsed
	}

	return &LoadBalancer{
		servers:      parsedServers,
		jwtValidator: jwtValidator,
	}
}

// Reverse proxy handler
func (lb *LoadBalancer) ReverseProxy(w http.ResponseWriter, r *http.Request) {
	claims, err := lb.jwtValidator.ValidateToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var targetServer *url.URL

	// Route admin requests to the first server only
	if claims.Role == "admin" {
		targetServer = lb.servers[0]
		log.Printf("Admin request detected - routing to admin server: %s", targetServer.String())
	} else {
		// For non-admin roles (client, user), use round-robin load balancing
		// but skip the first server (admin server)
		lb.mu.Lock()
		// Calculate index for non-admin servers (starting from index 1)
		if len(lb.servers) > 1 {
			// Use modulo on the remaining servers (excluding admin server)
			nonAdminIndex := (lb.currentIdx % (len(lb.servers) - 1)) + 1
			targetServer = lb.servers[nonAdminIndex]
			lb.currentIdx++
		} else {
			// Fallback if there's only one server
			targetServer = lb.servers[0]
		}
		lb.mu.Unlock()

		log.Printf("Routing %s request to %s", claims.Role, targetServer.String())
	}

	proxy := httputil.NewSingleHostReverseProxy(targetServer)
	proxy.ServeHTTP(w, r)
}
