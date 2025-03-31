package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gentrit-hoxha/golang-load-balancer/internal/jwt"
	"github.com/gentrit-hoxha/golang-load-balancer/internal/loadbalancer"

	"gopkg.in/yaml.v2"
)

// Configuration structure for loading YAML
type Config struct {
	LoadBalancer struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"load_balancer"`
	Servers []string `yaml:"servers"`
}

func main() {
	// Load configuration
	configFile, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Initialize JWT validator
	jwtValidator := jwt.NewValidator("3YGFHIEUFJUHEIFJEIFJM")

	// Create load balancer
	lb := loadbalancer.NewLoadBalancer(config.Servers, jwtValidator)

	// Start server
	log.Printf("Load Balancer listening on %s:%s", config.LoadBalancer.Host, config.LoadBalancer.Port)
	http.HandleFunc("/", lb.ReverseProxy)
	log.Fatal(http.ListenAndServe(":"+config.LoadBalancer.Port, nil))
}
