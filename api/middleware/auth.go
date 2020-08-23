package middleware

import (
	"context"
	"encoding/json"
	"github.com/nikolausreza131192/pos/controllers"
	"github.com/nikolausreza131192/pos/entity"
	"log"
	"net/http"
)

type AuthMW struct {
	AuthController controllers.Auth
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

const (
	userContextKey = "user"

	NotAuthorized      = "NOT AUTHORIZED"
	InvalidPermission  = "You don't have permission to access this resource"
	MissingAccessToken = "Missing access token"
)

func (mw *AuthMW) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response AuthResponse

		w.Header().Set("Content-Type", "application/json")
		token := r.Header.Get("x-access-token")
		if token == "" {
			w.WriteHeader(http.StatusForbidden)
			response.Success = false
			response.Message = MissingAccessToken
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
			return
		}

		userFromToken, err := mw.AuthController.GetUserFromToken(token)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			response.Success = false
			response.Message = NotAuthorized
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, userFromToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw *AuthMW) CheckPermission(f http.HandlerFunc, permissionName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response AuthResponse

		w.Header().Set("Content-Type", "application/json")
		// Get user from context
		userFromCtx := r.Context().Value(userContextKey)
		if userFromCtx == nil {
			w.WriteHeader(http.StatusForbidden)
			response.Success = false
			response.Message = NotAuthorized
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
			return
		}
		user, ok := userFromCtx.(entity.User)
		if !ok {
			log.Println("CheckPermission user from context is invalid", userFromCtx)
			w.WriteHeader(http.StatusForbidden)
			response.Success = false
			response.Message = NotAuthorized
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
			return
		}

		// Check permission
		isAuthorized := mw.AuthController.CheckPermissionAccess(user, permissionName)
		if !isAuthorized {
			w.WriteHeader(http.StatusForbidden)
			response.Success = false
			response.Message = InvalidPermission
			jsonResponse, _ := json.Marshal(response)
			w.Write(jsonResponse)
			return
		}

		f(w, r)
	}
}
