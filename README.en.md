# MantyBoot

A Go web framework inspired by Spring Boot, providing HTTP routing, middleware support, configuration management, and database error handling.

## Features

- **HTTP Router**: Compatible with Go's standard library
- **Middleware Support**: CORS, Rate Limiting, Recovery, Access Logging
- **Configuration Management**: YAML-based with validation
- **Structured Logging**: JSON and text formats with levels
- **Error Handling**: Comprehensive error types and wrapping
- **Database Support**: MySQL error translation

## Modules

| Module | Description | Install |
|--------|-------------|---------|
| [configuration](#configuration) | YAML configuration management like Spring Boot's ConfigurationProperties | `go get github.com/zbum/mantyboot/configuration` |
| [data](#data) | Database error abstraction and translation | `go get github.com/zbum/mantyboot/data` |
| [http](#http) | HTTP request handling utilities | `go get github.com/zbum/mantyboot/http` |

## Requirements

- Go 1.22 or higher

## Installation

```shell
go get github.com/zbum/mantyboot
```

## Quick Start

### Basic HTTP Server

```go
package main

import (
	"log"
	"net/http"

	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
)

func main() {
	mux := mux.NewMantyMux()

	// Add middleware
	mux.AddMiddleware(middleware.LegacyAccessLogger(log.Default()))

	// Register routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, MantyBoot!"))
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Enhanced Server with All Features

```go
package main

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/zbum/mantyboot/configuration"
	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
)

type AppConfig struct {
	Server struct {
		Port int    `yaml:"port" validate:"required,min=1,max=65535"`
		Host string `yaml:"host" validate:"required"`
	} `yaml:"server"`
}

//go:embed config/application-dev.yaml
var devfs embed.FS

func main() {
	// Initialize structured logger
	logger := log.Default()
	logger.Info("Starting enhanced application")

	// Load configuration with validation
	config, err := configuration.NewConfiguration[AppConfig](devfs, "dev")
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Create mux with enhanced middleware
	mux := mux.NewMantyMux()
	mux.AddMiddleware(middleware.Recovery(logger))
	mux.AddMiddleware(middleware.CORS(nil)) // Use default CORS config
	mux.AddMiddleware(middleware.RateLimitByIP(100, time.Minute))
	mux.AddMiddleware(middleware.AccessLogger(logger))

	// Register routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Enhanced MantyBoot!"))
	})

	// Start server
	addr := fmt.Sprintf("%s:%d", config.GetConfiguration().Server.Host, config.GetConfiguration().Server.Port)
	logger.Println("Starting enhanced application")

	logger.Fatal("Server failed", fmt.Errorf("server error: %v", http.ListenAndServe(addr, mux)))
}
```

---

## Configuration

YAML configuration management module similar to Spring Boot's `ConfigurationProperties`.

### File Structure
Configuration files are loaded in the following order. Later values override earlier ones.

1. Embedded files (embed.FS)
2. `./application-{profile}.yaml`
3. `./config/application-{profile}.yaml`

### Example Configuration

```yaml
# application-dev.yaml
server:
  port: 8080
  host: "localhost"

database:
  url: "mysql://localhost:3306/testdb"
  max-conns: 10
```

### Configuration with Validation

```go
type AppConfig struct {
	Server struct {
		Port int    `yaml:"port" validate:"required,min=1,max=65535"`
		Host string `yaml:"host" validate:"required"`
	} `yaml:"server"`
}

// Create validator
validator := configuration.NewConfigurationValidator()
validator.AddRule("Server.Port", configuration.ValidationRule{
	Required: true,
	Min:      &[]int{1}[0],
	Max:      &[]int{65535}[0],
})

// Load with validation
config, err := configuration.NewConfigurationWithValidation[AppConfig](devfs, "dev", validator)
```

---

## Data

Database error abstraction and translation module.

### MySQL Error Translator

Translates MySQL error codes into meaningful error types.

| Error Code | Error Type | Description |
|------------|------------|-------------|
| 1062 | `DuplicateKeyError` | Duplicate key error |
| 1452 | `FkConstraintError` | Foreign key constraint violation |

```go
import "github.com/zbum/mantyboot/data/mysql"

translator := mysql.MysqlErrorTranslator{}
translatedErr := translator.TranslateExceptionIfPossible(mysqlErr)

switch err := translatedErr.(type) {
case mysql.DuplicateKeyError:
	fmt.Printf("Duplicate key on table %s, column %s\n", err.Table, err.Column)
case mysql.FkConstraintError:
	fmt.Printf("Foreign key constraint failed on table %s\n", err.Table)
}
```

---

## HTTP

HTTP request handling utility module.

### RequestWrapper

Easily extract parameters from HTTP requests.

```go
package main

import (
    "net/http"
    mantyhttp "github.com/zbum/mantyboot/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    wrapper := mantyhttp.NewRequestWrapper(r)

    // Extract int64 parameter from Query String or POST Form
    id, err := wrapper.ParamInt64("id")
    if err != nil {
        // Handle error
    }

    // Other integer types are also supported
    count, _ := wrapper.ParamInt32("count")
    page, _ := wrapper.ParseInt("page")
}
```

### Parse Function

A generic function to convert strings to various integer types.

```go
import mantyhttp "github.com/zbum/mantyboot/http"

// Convert to int64
val, err := mantyhttp.Parse[int64]("12345")

// Convert to int32
val32, err := mantyhttp.Parse[int32]("123")
```

### MIME Type Constants

Provides commonly used Content-Type values as constants.

```go
import "github.com/zbum/mantyboot/http/mime"

contentType := mime.ContentTypeApplicationJson      // "application/json"
formType := mime.ContentTypeApplicationFormUrlencoded // "application/x-www-form-urlencoded"
```

---

## Middleware

### CORS
```go
corsConfig := middleware.DefaultCORSConfig()
corsConfig.AllowedOrigins = []string{"http://localhost:3000"}
mux.AddMiddleware(middleware.CORS(corsConfig))
```

### Rate Limiting
```go
// Limit by IP: 100 requests per minute
mux.AddMiddleware(middleware.RateLimitByIP(100, time.Minute))

// Custom rate limiting
limiter := middleware.NewRateLimiter(50, time.Hour)
mux.AddMiddleware(middleware.RateLimit(limiter, middleware.IPKeyFunc))
```

### Recovery
```go
// Basic recovery
mux.AddMiddleware(middleware.Recovery(logger))

// Custom recovery handler
mux.AddMiddleware(middleware.RecoveryWithHandler(logger, func(w http.ResponseWriter, r *http.Request, err interface{}) {
	http.Error(w, "Custom error message", http.StatusInternalServerError)
}))
```

### Access Logging
```go
logger := logging.Default()
mux.AddMiddleware(middleware.AccessLogger(logger))
```

---

## Error Handling

### Custom Error Types
```go
import "github.com/zbum/mantyboot/errors"

// Wrap errors with context
err := errors.WrapConfigurationError(originalErr, "failed to load config")

// Database errors
dbErr := errors.WrapDatabaseError(originalErr, "query", "failed to fetch user")

// HTTP errors
httpErr := errors.WrapHTTPError(originalErr, 500, "internal server error")
```

---

## Examples

See the `example/` directory for complete working examples:

- `example/http/main.go` - Basic HTTP server
- `example/configuration_ex1.go` - Configuration example
- `example/enhanced_demo.go` - Full-featured example with all components

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

Apache License 2.0

See [LICENSE](LICENSE) for details.