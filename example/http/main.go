package main

import (
	"log"
	"net/http"

	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
)

func main() {

	mux := mux.NewMantyMux()

	mux.AddMiddleware(middleware.AccessLogger(log.Default()))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	mux.HandleFunc("GET /manty", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("manty"))
	})

	log.Panic(http.ListenAndServe(":8080", mux))
}
