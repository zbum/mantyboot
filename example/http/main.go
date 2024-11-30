package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

type MantyMux struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewMantyMux() *MantyMux {
	return &MantyMux{mux: http.NewServeMux(), middlewares: []Middleware{}}
}

func (m *MantyMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mainHandler := handler
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		mainHandler = m.middlewares[i](mainHandler)
	}

	m.mux.HandleFunc(pattern, mainHandler)
}

func (m *MantyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func main() {

	mux := NewMantyMux()

	mux.middlewares = append(mux.middlewares, func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			fmt.Println(1, t)
			next(w, r)
		}
	})

	mux.middlewares = append(mux.middlewares, func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			fmt.Println(2, t)
			next(w, r)
		}
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	log.Panic(http.ListenAndServe(":8080", mux))
}
