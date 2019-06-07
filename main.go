package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	DBInit()
	defer DBClose()
	log.Printf("DB init ok")

	// setup server
	r := mux.NewRouter()
	r.Use(ResponseMiddleware)
	AllHandlerSetup(r)
	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// start background job
	go TaskGetWeatherData()

	log.Printf("start serve:")
	log.Fatal(server.ListenAndServe())
}
