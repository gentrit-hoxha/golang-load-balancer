package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims structure for JWT
type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func main() {
	// Define command line flags
	secret := flag.String("secret", "3YGFHIEUFJUHEIFJEIFJM", "Secret key for signing the token")
	expiry := flag.Int("expiry", 24, "Token expiry time in hours")
	host := flag.String("host", "localhost", "Host for curl commands")
	port := flag.String("port", "8080", "Port for curl commands")
	endpoint := flag.String("endpoint", "/api", "API endpoint for curl commands")
	flag.Parse()

	// Define roles
	roles := []string{"admin", "client", "user"}

	fmt.Println("=== JWT TOKENS ===")
	fmt.Println()

	// Create a file to save all tokens and curl commands
	f, err := os.Create("jwt_tokens.txt")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer f.Close()

	// Write header to file
	fmt.Fprintf(f, "JWT Tokens (expires in %d hours)\n", *expiry)
	fmt.Fprintf(f, "Generated on: %s\n\n", time.Now().Format(time.RFC1123))

	// Generate token for each role
	for _, role := range roles {
		// Create claims
		claims := &Claims{
			Role: role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(*expiry) * time.Hour).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign token
		tokenString, err := token.SignedString([]byte(*secret))
		if err != nil {
			log.Fatalf("Error signing token: %v", err)
		}

		// Create curl command
		curlCmd := fmt.Sprintf("curl -X GET http://%s:%s%s -H \"Authorization: Bearer %s\"",
			*host, *port, *endpoint, tokenString)

		// Output to console
		fmt.Printf("=== %s ===\n", role)
		fmt.Printf("Token: %s\n\n", tokenString)
		fmt.Printf("Curl command:\n%s\n\n", curlCmd)

		// Write to file
		fmt.Fprintf(f, "=== %s ===\n", role)
		fmt.Fprintf(f, "Token: %s\n\n", tokenString)
		fmt.Fprintf(f, "Curl command:\n%s\n\n", curlCmd)
	}

	fmt.Println("All tokens and curl commands saved to jwt_tokens.txt")
}
