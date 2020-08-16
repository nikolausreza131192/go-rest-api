package controllers_test

import (
	"errors"
	"testing"

	"github.com/nikolausreza131192/pos/controllers"
	"github.com/stretchr/testify/assert"

	"github.com/nikolausreza131192/pos/repository"
)

func TestCreateUser(t *testing.T) {
	tcs := []struct {
		name                   string
		nameParam              string
		usernameParam          string
		emailParam             string
		roleParam              string
		userRepo               repository.UserRepo
		expectedPasswordResult string
		expectedError          error
	}{
		{
			name:          "Invalid parameter: empty name",
			userRepo:      &fakeUserRepo{},
			expectedError: errors.New("Invalid parameter"),
		},
		{
			name:          "Error on creating user",
			nameParam:     "Foo Bar",
			usernameParam: "foobar",
			emailParam:    "foo@bar.com",
			roleParam:     "admin",
			userRepo: &fakeUserRepo{
				CreateUserError: errors.New("Some error"),
			},
			expectedError: errors.New("Some error"),
		},
		{
			name:          "Successfully create user",
			nameParam:     "Foo Bar",
			usernameParam: "foobar",
			emailParam:    "foo@bar.com",
			roleParam:     "admin",
			userRepo: &fakeUserRepo{
				CreateUserResult: "hashedrandompassword",
			},
			expectedPasswordResult: "hashedrandompassword",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			param := controllers.UserControllerParam{
				UserRP: tc.userRepo,
			}
			controller := controllers.NewUser(param)
			password, err := controller.CreateUser(tc.nameParam, tc.usernameParam, tc.emailParam, tc.roleParam)

			assert.Equal(t, tc.expectedPasswordResult, password)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
