# 🔄 Golang Load Balancer

A load balancer implementation in Go that distributes incoming HTTP requests across multiple backend servers.

## 🚀 Project Overview

This project implements a load balancer in Go that can:
- 🔀 Distribute traffic across multiple backend servers
- 🔐 Handle authentication with JWT tokens
- 📊 Provide basic monitoring and health checks
- 🔥 Support load testing

## 🛠️ Setup Instructions

### Prerequisites
- Go 1.16+
- Python 3.6+ (for testing scripts)
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/golang-load-balancer.git
   cd golang-load-balancer
   ```

2. Set up Python environment (for testing):
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   # Install required Python packages if needed
   # pip install -r requirements.txt
   ```

3. Build the load balancer:
   ```bash
   go build -o load-balancer cmd/load-balancer/main.go
   ```

## 📋 Usage

### Starting Backend Servers

Start multiple backend test servers on different ports:

```bash
# Start server on port 8081
go run scripts/start_server/test_server.go 8081

# Start server on port 8082
go run scripts/start_server/test_server.go 8082

# Start server on port 8083
go run scripts/start_server/test_server.go 8083
```

### Running the Load Balancer

Start the load balancer:

```bash
go run cmd/load-balancer/main.go
```

Or use the compiled binary:

```bash
./load-balancer
```

### 🔑 Authentication

Generate JWT authentication tokens:

```bash
go run generate_token.go
```

### 🔍 Load Testing

Run the load testing script to simulate multiple requests:

```bash
./scripts/load_test.sh 100
```

## ⚙️ Configuration

The load balancer can be configured by modifying the configuration file (default: `configs/config.yaml`). You can specify:
- Backend server addresses
- Load balancing algorithm (round-robin, least connections, etc.)
- Health check intervals
- JWT authentication settings

## 📁 Project Structure

```
jwt-load-balancer/
│
├── cmd/
│   └── load-balancer/
│       └── main.go             # Entry point for load balancer
│
├── internal/
│   ├── loadbalancer/
│   │   └── loadbalancer.go     # Core load balancing logic
│   └── jwt/
│       └── jwt.go              # JWT authentication and validation
│
├── configs/
│   └── config.yaml             # Project configuration
│
├── scripts/
│   ├── load_test.sh            # Load testing bash script
│   └── analyze_logs.py         # Log analysis Python script
│
├── tests/
│   └── loadbalancer_test.go    # Unit tests
│
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency lockfile
└── README.md                   # Project documentation
```

## 🔒 Security Features

- JWT-based authentication system
- Token validation and expiration handling
- Request authentication

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.