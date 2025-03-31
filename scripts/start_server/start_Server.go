package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Configure logger with custom format
	log.SetFlags(0) // Remove default timestamp

	// Check if port is provided as command line argument
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run start_Server.go [port]")
	}

	// Parse port from command line
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Create a unique server identifier
	serverID := fmt.Sprintf("Server-%d", port)

	// Set different colors for different servers
	var colorCode string
	switch port {
	case 8081:
		colorCode = "\033[1;31m" // Red for admin server
	case 8082:
		colorCode = "\033[1;32m" // Green for second server
	case 8083:
		colorCode = "\033[1;34m" // Blue for third server
	default:
		colorCode = "\033[1;35m" // Magenta for other servers
	}
	resetColor := "\033[0m"

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Define HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Record start time
		startTime := time.Now()

		// Extract role from Authorization header if present
		role := "unknown"
		authHeader := r.Header.Get("Authorization")
		if strings.Contains(authHeader, "eyJ") { // Simple check for JWT format
			role = extractRoleFromHeader(authHeader)
		}

		// Simulate processing time (between 50-200ms)
		processingTime := 50 + rand.Intn(150)
		time.Sleep(time.Duration(processingTime) * time.Millisecond)

		// Simulate database action

		// Simulate database time (between 10-100ms)
		dbTime := 10 + rand.Intn(90)
		time.Sleep(time.Duration(dbTime) * time.Millisecond)

		// Calculate total processing time
		totalTime := time.Since(startTime).Milliseconds()

		// Log the request in a single, informative line with fixed width columns for better readability
		log.Printf("%s[%s]%s %-6s | %-15s | Role: %-7s | Process: %3dms | DB: %3dms | Total: %3dms",
			colorCode, serverID, resetColor,
			r.Method,
			r.RemoteAddr,
			role,
			processingTime,
			dbTime,
			totalTime)

		// Return response with server information
		fmt.Fprintf(w, "Response from %s\n", serverID)
		fmt.Fprintf(w, "Time: %s\n", time.Now().Format(time.RFC1123))
		fmt.Fprintf(w, "Request path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Remote address: %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Role: %s\n", role)
		fmt.Fprintf(w, "Processing time: %dms\n", totalTime)
	})

	// Start the server
	serverAddr := fmt.Sprintf(":%d", port)

	// Print server startup message
	log.Printf("%s[%s]%s Starting server on http://127.0.0.1%s",
		colorCode, serverID, resetColor, serverAddr)
	log.Printf("%s[%s]%s Server role: %s",
		colorCode, serverID, resetColor, getServerRole(port))
	log.Printf("%s[%s]%s Press Ctrl+C to stop",
		colorCode, serverID, resetColor)

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("%s[%s]%s Server error: %v",
			colorCode, serverID, resetColor, err)
	}
}

// Simple function to extract role from JWT token
func extractRoleFromHeader(authHeader string) string {
	// Remove "Bearer " if present
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Split the token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "invalid-token"
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "invalid-payload"
	}

	// Parse the JSON payload
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "invalid-json"
	}

	// Check for role in claims
	if role, ok := claims["role"].(string); ok {
		return role
	}

	return "no-role"
}

// Get server role based on port
func getServerRole(port int) string {
	switch port {
	case 8081:
		return "Admin Server (handles admin, client, and user requests)"
	case 8082, 8083:
		return "Regular Server (handles client and user requests)"
	default:
		return "Unknown Server Role"
	}
}

// Helper function to truncate strings for display
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}
