package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/zbum/mantyboot/http/mux"
)

func Recovery(logger *log.Logger) mux.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic
					logger.Println("panic recovered", fmt.Errorf("panic: %v", err))
					logger.Println("stack trace: " + string(debug.Stack()))

					// Return 500 Internal Server Error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next(w, r)
		}
	}
}

func RecoveryWithHandler(logger *log.Logger, handler func(http.ResponseWriter, *http.Request, interface{})) mux.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic
					logger.Println("panic recovered", fmt.Errorf("panic: %v", err))
					logger.Println("stack trace: " + string(debug.Stack()))

					// Call custom handler
					handler(w, r, err)
				}
			}()

			next(w, r)
		}
	}
}
