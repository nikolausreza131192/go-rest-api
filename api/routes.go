package main

import (
	"fmt"
	"net/http"

	"github.com/nikolausreza131192/pos/api/handler"

	"github.com/gorilla/mux"
	"github.com/nikolausreza131192/pos/entity"
)

// GetAllItemsResponse is main response for get all items
type GetAllItemsResponse struct {
	Data []entity.Item `json:"data"`
}

func initRoutes(r *mux.Router, controllers Controllers) {
	fmt.Println("Init routes...")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`POS API Ready`))
	})
	subRouter := r.PathPrefix("/api/v1").Subrouter()
	authenticatedRoutes := subRouter
	authenticatedRoutes.HandleFunc("/items", handler.GetAllItems(controllers.Item)).Methods("GET")
	authenticatedRoutes.HandleFunc("/items/{id}", handler.GetItemByID(controllers.Item)).Methods("GET")
	authenticatedRoutes.HandleFunc("/user", handler.CreateUser(controllers.User)).Methods("POST")
}
