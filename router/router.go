package router

import (
	"first-go-postgres/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user", middleware.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.DeleteUsers).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")

	return router
}
