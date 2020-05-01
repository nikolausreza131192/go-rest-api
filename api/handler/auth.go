package handler

import (
	"net/http"

	"github.com/nikolausreza131192/pos/controllers"
)

// GetToken return user's token
func GetToken(controller controllers.Item) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		return
	}
}
