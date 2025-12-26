package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zbum/mantyboot/configuration"
	"github.com/zbum/mantyboot/errors"
	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
)

// Enhanced configuration with validation
type EnhancedConfiguration struct {
	Server struct {
		Port int    `yaml:"port" validate:"required,min=1,max=65535"`
		Host string `yaml:"host" validate:"required"`
	} `yaml:"server"`
	Database struct {
		URL      string `yaml:"url" validate:"required"`
		MaxConns int    `yaml:"max-conns" validate:"min=1,max=100"`
	} `yaml:"database"`
}

//go:embed config/application-dev.yaml
var devfs embed.FS

func main() {
	// Initialize structured logger
	logger := log.Default()
	logger.Println("Starting enhanced application")

	// Load configuration with validation
	config, err := loadConfiguration()
	if err != nil {
		logger.Fatal("Failed to load configuration", err)
	}

	// Create mux with enhanced middleware
	mux := createEnhancedMux(logger)

	// Register routes
	registerRoutes(mux, logger)

	// Start server
	addr := fmt.Sprintf("%s:%d", config.GetConfiguration().Server.Host, config.GetConfiguration().Server.Port)

	logger.Fatal("Server failed", fmt.Errorf("server error: %v", http.ListenAndServe(addr, mux)))
}

func loadConfiguration() (*configuration.Configuration[EnhancedConfiguration], error) {
	// Create validator with custom rules
	validator := configuration.NewConfigurationValidator()

	// Add custom validation rules
	validator.AddRule("Server.Port", configuration.ValidationRule{
		Required: true,
		Min:      &[]int{1}[0],
		Max:      &[]int{65535}[0],
		Custom: func(value interface{}) error {
			if port, ok := value.(int); ok {
				return configuration.ValidatePort(port)
			}
			return fmt.Errorf("invalid port type")
		},
	})

	validator.AddRule("Database.URL", configuration.ValidationRule{
		Required: true,
		Custom: func(value interface{}) error {
			if url, ok := value.(string); ok {
				return configuration.ValidateURL(url)
			}
			return fmt.Errorf("invalid URL type")
		},
	})

	// Load configuration with validation
	config, err := configuration.NewConfigurationWithValidation[EnhancedConfiguration](devfs, "dev", validator)
	if err != nil {
		return nil, errors.WrapConfigurationError(err, "failed to load configuration")
	}

	return config, nil
}

func createEnhancedMux(logger *log.Logger) *mux.MantyMux {
	mux := mux.NewMantyMux()

	// Add recovery middleware
	mux.AddMiddleware(middleware.Recovery(logger))

	// Add CORS middleware
	corsConfig := middleware.DefaultCORSConfig()
	corsConfig.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	mux.AddMiddleware(middleware.CORS(corsConfig))

	// Add rate limiting middleware
	mux.AddMiddleware(middleware.RateLimitByIP(100, time.Minute))

	// Add structured access logging
	mux.AddMiddleware(middleware.AccessLogger(logger))

	return mux
}

func registerRoutes(mux *mux.MantyMux, logger *log.Logger) {
	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	// API endpoints
	mux.HandleFunc("GET /api/users", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]`))
	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {

		// Simulate some processing
		time.Sleep(100 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":3,"name":"New User"}`))
	})

	// Error simulation endpoint
	mux.HandleFunc("GET /api/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	// Panic simulation endpoint
	mux.HandleFunc("GET /api/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("simulated panic")
	})

	// Root endpoint
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head><title>Enhanced MantyBoot Example</title></head>
			<body>
				<h1>Enhanced MantyBoot Example</h1>
				<p>This example demonstrates:</p>
				<ul>
					<li>Structured logging</li>
					<li>Configuration validation</li>
					<li>CORS middleware</li>
					<li>Rate limiting</li>
					<li>Recovery middleware</li>
					<li>Enhanced error handling</li>
				</ul>
				<p><a href="/health">Health Check</a></p>
				<p><a href="/api/users">Users API</a></p>
			</body>
			</html>
		`))
	})
}
