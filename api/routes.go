package main

import (
	"fmt"
	"net/http"

	"github.com/nikolausreza131192/pos/api/handler"

	"github.com/gorilla/mux"
)

func initRoutes(r *mux.Router, controllers Controllers) {
	fmt.Println("Init routes...")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`POS API Ready`))
	})
	subRouter := r.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/login", handler.Login(controllers.Auth)).Methods("POST")

	authenticatedRoutes := subRouter
	authenticatedRoutes.HandleFunc("/items", handler.GetAllItems(controllers.Item)).Methods("GET")
	authenticatedRoutes.HandleFunc("/items/{id}", handler.GetItemByID(controllers.Item)).Methods("GET")
	authenticatedRoutes.HandleFunc("/user", handler.CreateUser(controllers.User)).Methods("POST")
}
