# Golang Library - Reusable Components

A collection of reusable Go packages for building production-ready applications.

## 📦 Packages

### Logger
Structured logging with Zap integration, supporting multiple environments and log levels.

**Features:**
- Structured logging with context support
- Configurable log levels (debug, info, warn, error)
- Environment-specific configurations (development, production)
- Request ID tracking
- Service tagging

**Usage:**
```go
import "github.com/bikashb-meesho/golang-lib/logger"

// Create a logger
log, err := logger.New(logger.Config{
    Level:       "info",
    Environment: "production",
    Service:     "my-service",
})

log.Info("User created", zap.String("user_id", "123"))
```

### Validator
Comprehensive input validation with common validation rules.

**Features:**
- Required field validation
- Length constraints (min/max)
- Email validation
- Range validation
- Pattern matching (regex)
- Custom validation rules

**Usage:**
```go
import "github.com/bikashb-meesho/golang-lib/validator"

v := validator.New()
v.Required("name", user.Name)
v.Email("email", user.Email)
v.MinLength("password", user.Password, 8)

if !v.IsValid() {
    return fmt.Errorf("validation errors: %s", v.ErrorMessages())
}
```

### HTTP Utilities
HTTP helper functions and middleware for building REST APIs.

**Features:**
- JSON response helpers
- Request body parsing with size limits
- Recovery middleware (panic handling)
- Request ID middleware
- Timeout middleware
- CORS middleware

**Usage:**
```go
import "github.com/bikashb-meesho/golang-lib/httputil"

func handler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := httputil.ReadJSON(r, &req, 1<<20); err != nil {
        httputil.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
        return
    }
    
    httputil.WriteSuccess(w, map[string]string{"id": "123"})
}

// With middleware
mux := http.NewServeMux()
mux.HandleFunc("/api/users", handler)

handler := httputil.Recover(
    httputil.RequestID(
        httputil.Timeout(30 * time.Second)(mux),
    ),
)
```

## 🚀 Installation

```bash
go get github.com/bikashb-meesho/golang-lib
```

## 🔧 Development

### Prerequisites
- Go 1.23 or higher

### Running Tests
```bash
go test ./...
```

### Running Tests with Coverage
```bash
go test -cover ./...
```

## 📝 GitHub Setup

### Creating the Repository

1. Create a new GitHub repository named `golang-lib`

2. Initialize git in this directory:
```bash
cd golang-lib
git init
git add .
git commit -m "Initial commit: Add reusable library components"
```

3. Add remote and push:
```bash
git remote add origin https://github.com/bikashb-meesho/golang-lib.git
git branch -M main
git push -u origin main
```

4. Tag a release:
```bash
git tag v1.0.0
git push origin v1.0.0
```

## 📄 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

