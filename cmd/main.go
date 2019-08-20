package main

import (
	"github.com/gorilla/mux"
	"github.com/gyanesh-m/File-download-manager/src/route"
	"log"
	"net/http"
	"time"
)


func main() {
	router := mux.NewRouter().StrictSlash(true)
	route.HandleRequests(router)
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8081",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}