package route

import (
	"../controller"
	"github.com/gorilla/mux"
)

func HandleRequests(router *mux.Router) {
	router.HandleFunc("/health", controller.HomePage).Methods("GET")
	// error handing in get and post for downloads.
	router.HandleFunc("/downloads/{id}", controller.Status).Methods("GET")
	router.HandleFunc("/downloads", controller.Download).Methods("POST")
	// browse files
	router.HandleFunc("/files", controller.Browse).Methods("GET")
}
