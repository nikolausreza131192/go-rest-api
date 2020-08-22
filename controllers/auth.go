package controllers

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// User control all process related with user
type Auth interface {
	Login(username, password string) (entity.LoginResponse, error)
}

type authController struct {
	userRP    repository.UserRepo
	authRP    repository.AuthRepo
	loginTime int // in minute
}

// UserControllerParam will be used as repository parameter
type AuthControllerParam struct {
	UserRP    repository.UserRepo
	AuthRP    repository.AuthRepo
	LoginTime int // in minute
}

// NewUser initialize user controller
func NewAuth(param AuthControllerParam) Auth {
	return &authController{
		userRP:    param.UserRP,
		authRP:    param.AuthRP,
		loginTime: param.LoginTime,
	}
}

func (c *authController) Login(username, password string) (entity.LoginResponse, error) {
	var response entity.LoginResponse

	// Validate param
	if username == "" {
		response = entity.LoginResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Username can not be empty",
		}
		return response, errors.New("username can not be empty")
	}
	if password == "" {
		response = entity.LoginResponse{
			Success: false,
			Code:    http.StatusBadRequest,
			Message: "Password can not be empty",
		}
		return response, errors.New("password can not be empty")
	}

	// Find user
	user, err := c.userRP.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			response = entity.LoginResponse{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Invalid username/password",
			}
			return response, nil
		}
		log.Println("func Login Error get user", username, err)
		response = entity.LoginResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while processing your request.",
		}
		return response, err
	}

	// Get the password
	userPassword, err := c.userRP.GetUserPassword(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			response = entity.LoginResponse{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Invalid username/password",
			}
			return response, nil
		}
		log.Println("func Login Error get password", user.Username, err)
		response = entity.LoginResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while processing your request.",
		}
		return response, err
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			response = entity.LoginResponse{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Invalid username/password",
			}
			return response, nil
		}
		log.Println("func Login Error compare password", user.Username, err)
		response = entity.LoginResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while processing your request.",
		}
		return response, err
	}

	// Get JWT token
	token, err := c.authRP.GenerateJWTToken(entity.TokenClaims{
		Name:     user.Name,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(c.loginTime) * time.Minute).Unix(),
		},
	})
	if err != nil {
		log.Println("func Login Error generate token", user.Username, err)
		response = entity.LoginResponse{
			Success: false,
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while processing your request.",
		}
		return response, err
	}

	response = entity.LoginResponse{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Token:   token,
		User:    user,
	}

	return response, nil
}
