package controllers

import (
	"errors"
	"log"
	"time"

	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
)

// User control all process related with user
type User interface {
	CreateUser(name, username, email, role string) (string, error)
}

type userController struct {
	userRP repository.UserRepo
}

// UserControllerParam will be used as repository parameter
type UserControllerParam struct {
	UserRP repository.UserRepo
}

// NewUser initialize user controller
func NewUser(param UserControllerParam) User {
	return &userController{
		userRP: param.UserRP,
	}
}

func (c *userController) CreateUser(name, username, email, role string) (string, error) {
	// Validate param first
	if name == "" || username == "" || email == "" || role == "" {
		log.Printf("func CreateUser Invalid parameter. Name: %s. Username: %s. Email: %s. Role: %s", name, username, email, role)
		return "", errors.New("Invalid parameter")
	}

	now := time.Now()
	userData := entity.User{
		Name:      name,
		Username:  username,
		Email:     email,
		Role:      role,
		CreatedBy: "test",
		UpdatedBy: "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	password, err := c.userRP.CreateUser(userData)
	if err != nil {
		log.Printf("func CreateUser Error on creating user. User: %+v. Error: %s", userData, err)
		return "", err
	}

	return password, nil
}
