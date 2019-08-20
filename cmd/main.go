package main

import (
	"../src/route"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)


func main() {
	router := mux.NewRouter().StrictSlash(true)
	route.HandleRequests(router)
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}