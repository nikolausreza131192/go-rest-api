package entity

import "github.com/dgrijalva/jwt-go"

type (
	LoginResponse struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"` // represent http status code
		Message string `json:"message"`
		Token   string `json:"token"`
		User    `json:"user"`
	}

	TokenClaims struct {
		Name     string
		Username string
		jwt.StandardClaims
	}
)
