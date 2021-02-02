package router

import (
	"go-postgres-practice/middleware/user"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go to define and handle all endpoints and their respective middleware
func Router() *mux.Router {
	var router *mux.Router

	router = mux.NewRouter()

	router.HandleFunc("/api/user/{id}", user.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user", user.GetAllUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newuser", user.CreateUser).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{id}", user.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteuser/{id}", user.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}
