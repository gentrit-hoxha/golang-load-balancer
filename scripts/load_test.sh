#!/bin/bash

# Configuration
LOAD_BALANCER_URL="http://localhost:8080"
TOTAL_REQUESTS=${1:-100}
LOG_FILE="load_test.log"

# Clear previous log file
> "$LOG_FILE"

echo "Generating JWT tokens..."
# Generate tokens using the Go script
go run scripts/generate_token.go > tokens.txt

# Extract tokens from the output file
ADMIN_TOKEN=$(grep -A 1 "=== admin ===" tokens.txt | grep "Token:" | cut -d' ' -f2)
CLIENT_TOKEN=$(grep -A 1 "=== client ===" tokens.txt | grep "Token:" | cut -d' ' -f2)
USER_TOKEN=$(grep -A 1 "=== user ===" tokens.txt | grep "Token:" | cut -d' ' -f2)

if [ -z "$ADMIN_TOKEN" ] || [ -z "$CLIENT_TOKEN" ] || [ -z "$USER_TOKEN" ]; then
    echo "Error: Failed to extract tokens. Check if the token generation script is working."
    exit 1
fi

echo "Starting load test with $TOTAL_REQUESTS mixed requests..."

# Function to send a request with a specific role
send_request() {
    local role=$1
    local token=$2
    
    # Send request and get status code
    local status=$(curl -s -o /dev/null -w "%{http_code}" \
        -H "Authorization: Bearer $token" \
        "$LOAD_BALANCER_URL")
    
    # Log the result with role prefix
    echo "${role}:${status}" >> "$LOG_FILE"
}

# Send mixed requests
for ((i=1; i<=TOTAL_REQUESTS; i++)); do
    # Randomly select a role (weighted distribution: 20% admin, 40% client, 40% user)
    RANDOM_NUM=$((RANDOM % 100))
    
    if [ $RANDOM_NUM -lt 20 ]; then
        # Admin request (20% chance)
        send_request "admin" "$ADMIN_TOKEN" &
    elif [ $RANDOM_NUM -lt 60 ]; then
        # Client request (40% chance)
        send_request "client" "$CLIENT_TOKEN" &
    else
        # User request (40% chance)
        send_request "user" "$USER_TOKEN" &
    fi
    
    # Add a small delay between starting requests (50ms)
    sleep 0.07
    
    # Show progress every 10 requests
    if [ $((i % 10)) -eq 0 ] || [ $i -eq $TOTAL_REQUESTS ]; then
        echo "Progress: $i/$TOTAL_REQUESTS requests sent"
    fi
done

# Wait for all background processes to complete
echo "Waiting for all requests to complete..."
wait

echo "Load test completed. Results saved to $LOG_FILE"
echo "Analyzing results..."

# Run the Python analysis script
python3 scripts/analyze_logs.py "$LOG_FILE"