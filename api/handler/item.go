package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nikolausreza131192/pos/controllers"

	"github.com/nikolausreza131192/pos/entity"
)

// GetAllItemsResponse is main response for get all items
type GetAllItemsResponse struct {
	Data []entity.Item `json:"data"`
}

// GetAllItems is API handler to get all products
func GetAllItems(controller controllers.Item) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		items := controller.GetAll()
		response := GetAllItemsResponse{
			Data: items,
		}
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		return
	}
}

// GetItemByIDResponse is main response for get all items
type GetItemByIDResponse struct {
	Data  entity.Item `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// GetItemByID is API handler to get item by its ID
func GetItemByID(controller controllers.Item) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := GetItemByIDResponse{}
		defer func() {
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
		}()
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		itemID, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println("Handler: GetItemByID Can't convert URL param to int", vars["id"], err)
			w.WriteHeader(http.StatusBadRequest)
			response.Error = errors.New("Invalid Item ID").Error()
			return
		}
		item := controller.GetByID(itemID)
		response.Data = item
		w.WriteHeader(http.StatusOK)

		return
	}
}
