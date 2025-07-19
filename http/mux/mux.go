package mux

import (
	"log"
	"net/http"
	"os"

	"github.com/zbum/mantyboot/utils"
)

type MantyMux struct {
	mux         *http.ServeMux
	middlewares []Middleware
	logger      log.Logger
}

func NewMantyMux() *MantyMux {
	return &MantyMux{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
		logger:      *log.New(os.Stderr, "[mantyboot]", log.LstdFlags),
	}
}

func (m *MantyMux) AddMiddleware(middleware Middleware) {
	m.middlewares = append(m.middlewares, middleware)
}

func (m *MantyMux) Handle(pattern string, handler http.Handler) {

	mainHandlerFunc := handler.ServeHTTP
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		mainHandlerFunc = m.middlewares[i](mainHandlerFunc)
	}
	m.logger.Printf("[handler registered] pattern \"%s\", %s \n", pattern, utils.GetFunctionName(handler))

	m.mux.HandleFunc(pattern, mainHandlerFunc)
}

func (m *MantyMux) HandleFunc(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	mainHandler := handlerFunc
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		mainHandler = m.middlewares[i](mainHandler)
	}
	m.logger.Printf("[handlerFunc registered] pattern \"%s\", %s \n", pattern, utils.GetFunctionName(handlerFunc))

	m.mux.HandleFunc(pattern, mainHandler)
}

func (m *MantyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}
