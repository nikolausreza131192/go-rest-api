package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nikolausreza131192/pos/controllers"
)

// CreateUserResponse is main response for creating user API
type CreateUserResponse struct {
	UserPassword string `json:"user_password,omitempty"`
	Error        string `json:"error"`
}

// CreateUserRequest define request structure for creating user
type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// CreateUser is API handler to create new user
func CreateUser(controller controllers.User) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := CreateUserResponse{}
		defer func() {
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
		}()

		request := CreateUserRequest{}
		json.NewDecoder(r.Body).Decode(&request)
		w.Header().Set("Content-Type", "application/json")
		userPassword, err := controller.CreateUser(request.Name, request.Username, request.Email, request.Role)
		if err != nil {
			log.Printf("func CreateUser Error on creating user. Request: %+v. Error: %s", request, err)
			w.WriteHeader(http.StatusBadRequest)
			response.Error = err.Error()
			return
		}
		response.UserPassword = userPassword

		w.WriteHeader(http.StatusOK)
		return
	}
}
