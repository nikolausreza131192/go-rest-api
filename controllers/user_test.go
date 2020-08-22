package controllers_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/nikolausreza131192/pos/controllers"
	"github.com/stretchr/testify/assert"

	"github.com/nikolausreza131192/pos/repository"
)

func TestCreateUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tcs := []struct {
		name                   string
		nameParam              string
		usernameParam          string
		emailParam             string
		roleParam              string
		mockUserRepo           func() repository.UserRepo
		expectedPasswordResult string
		expectedError          error
	}{
		{
			name: "Invalid parameter: empty name",
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				return userRepo
			},
			expectedError: errors.New("Invalid parameter"),
		},
		{
			name:          "Error on creating user",
			nameParam:     "Foo Bar",
			usernameParam: "foobar",
			emailParam:    "foo@bar.com",
			roleParam:     "admin",
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				userRepo.EXPECT().CreateUser(gomock.Any()).Return("", errors.New("test error"))

				return userRepo
			},
			expectedError: errors.New("test error"),
		},
		{
			name:          "Successfully create user",
			nameParam:     "Foo Bar",
			usernameParam: "foobar",
			emailParam:    "foo@bar.com",
			roleParam:     "admin",
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				userRepo.EXPECT().CreateUser(gomock.Any()).Return("hashedrandompassword", nil)

				return userRepo
			},
			expectedPasswordResult: "hashedrandompassword",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			param := controllers.UserControllerParam{
				UserRP: tc.mockUserRepo(),
			}
			controller := controllers.NewUser(param)
			password, err := controller.CreateUser(tc.nameParam, tc.usernameParam, tc.emailParam, tc.roleParam)

			assert.Equal(t, tc.expectedPasswordResult, password)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
