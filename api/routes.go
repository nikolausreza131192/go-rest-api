package main

import (
	"fmt"
	"github.com/nikolausreza131192/pos/api/middleware"
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

	// Define auth middleware
	authMW := middleware.AuthMW{
		AuthController: controllers.Auth,
	}

	// Register authenticated routes
	authenticatedRoutes := r.PathPrefix("/api/v1").Subrouter()
	authenticatedRoutes.Use(authMW.Auth)
	authenticatedRoutes.HandleFunc("/items", authMW.CheckPermission(handler.GetAllItems(controllers.Item), "get list items")).Methods("GET")
	authenticatedRoutes.HandleFunc("/items/{id}", authMW.CheckPermission(handler.GetItemByID(controllers.Item), "get item details")).Methods("GET")
	authenticatedRoutes.HandleFunc("/user", handler.CreateUser(controllers.User)).Methods("POST")
}
