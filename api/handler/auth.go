package handler

import (
	"encoding/json"
	"github.com/nikolausreza131192/pos/entity"
	"log"
	"net/http"

	"github.com/nikolausreza131192/pos/controllers"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success     bool   `json:"success,omitempty"`
	Message     string `json:"message,omitempty"`
	Token       string `json:"token,omitempty"`
	entity.User `json:"user,omitempty"`
}

func Login(controller controllers.Auth) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := LoginResponse{}
		defer func() {
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
		}()

		request := LoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Println("func Login error parse JSON body", err)
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = "Request can not be empty"
			return
		}

		loginResponse, err := controller.Login(request.Username, request.Password)
		if err != nil {
			log.Println("func Login error from controller", request.Username, err)
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = err.Error()
			return
		}
		if !loginResponse.Success {
			w.WriteHeader(loginResponse.Code)
			response.Message = loginResponse.Message
			return
		}
		response.Success = loginResponse.Success
		response.Message = loginResponse.Message
		response.Token = loginResponse.Token
		response.User = loginResponse.User
		w.WriteHeader(loginResponse.Code)

		return
	}
}
