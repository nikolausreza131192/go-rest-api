package controllers_test

import (
	"database/sql"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/nikolausreza131192/pos/controllers"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tcs := []struct {
		name           string
		username       string
		password       string
		mockAuthRepo   func() repository.AuthRepo
		mockUserRepo   func() repository.UserRepo
		expectedResult entity.LoginResponse
		expectError    bool
	}{
		// Username is empty
		{
			name:     "Username is empty",
			username: "",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusBadRequest,
				Message: "Username can not be empty",
			},
			expectError: true,
		},
		// Password is empty
		{
			name:     "Password is empty",
			username: "foo",
			password: "",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusBadRequest,
				Message: "Password can not be empty",
			},
			expectError: true,
		},
		// Error on get user
		{
			name:     "Error on get user",
			username: "foo",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				userRepo.EXPECT().GetByUsername("foo").Return(entity.User{}, errors.New("test error"))

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusInternalServerError,
				Message: "An error occurred while processing your request.",
			},
			expectError: true,
		},
		// User not found
		{
			name:     "User not found",
			username: "foo",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				userRepo.EXPECT().GetByUsername("foo").Return(entity.User{}, sql.ErrNoRows)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Invalid username/password",
			},
			expectError: false,
		},
		// Error on get password
		{
			name:     "Error on get password",
			username: "foo",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				user := entity.User{
					Username: "foo",
					Name:     "Foo Bar",
				}
				userRepo.EXPECT().GetByUsername("foo").Return(user, nil)
				userRepo.EXPECT().GetUserPassword("foo").Return("", errors.New("test error"))

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusInternalServerError,
				Message: "An error occurred while processing your request.",
			},
			expectError: true,
		},
		// User not found on get password
		{
			name:     "User not found on get password",
			username: "foo",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				user := entity.User{
					Username: "foo",
					Name:     "Foo Bar",
				}
				userRepo.EXPECT().GetByUsername("foo").Return(user, nil)
				userRepo.EXPECT().GetUserPassword("foo").Return("", sql.ErrNoRows)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Invalid username/password",
			},
			expectError: false,
		},
		// Error validate password
		{
			name:     "Error validate password",
			username: "foobar",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				user := entity.User{
					Username: "foobar",
					Name:     "Foo Bar",
				}
				userRepo.EXPECT().GetByUsername("foobar").Return(user, nil)
				userRepo.EXPECT().GetUserPassword("foobar").Return("", errors.New("test error"))

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusInternalServerError,
				Message: "An error occurred while processing your request.",
			},
			expectError: true,
		},
		// Wrong password
		{
			name:     "Wrong password",
			username: "foo",
			password: "bar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				user := entity.User{
					Username: "foo",
					Name:     "Foo Bar",
				}
				userRepo.EXPECT().GetByUsername("foo").Return(user, nil)
				userRepo.EXPECT().GetUserPassword("foo").Return("$2a$10$JjIFwYYQ2ZiMe.dM1Dd0FuDm/2285KQ6Fk1xJhX9QWrGtOcT531Ty", nil)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusForbidden,
				Message: "Invalid username/password",
			},
			expectError: false,
		},
		// Error on generate token
		{
			name:     "Error on generate token",
			username: "foo",
			password: "foobar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				authRepo.EXPECT().GenerateJWTToken(gomock.Any()).Return("", errors.New("test error"))

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				user := entity.User{
					Username: "foo",
					Name:     "Foo Bar",
				}
				userRepo.EXPECT().GetByUsername("foo").Return(user, nil)
				userRepo.EXPECT().GetUserPassword("foo").Return("$2a$10$JjIFwYYQ2ZiMe.dM1Dd0FuDm/2285KQ6Fk1xJhX9QWrGtOcT531Ty", nil)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: false,
				Code:    http.StatusInternalServerError,
				Message: "An error occurred while processing your request.",
			},
			expectError: true,
		},
		// Successfully login generate token
		{
			name:     "Successfully login and generate token",
			username: "foo",
			password: "foobar",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				authRepo.EXPECT().GenerateJWTToken(gomock.Any()).Return("publictoken", nil)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				user := entity.User{
					Username: "foo",
					Name:     "Foo Bar",
				}
				userRepo.EXPECT().GetByUsername("foo").Return(user, nil)
				userRepo.EXPECT().GetUserPassword("foo").Return("$2a$10$JjIFwYYQ2ZiMe.dM1Dd0FuDm/2285KQ6Fk1xJhX9QWrGtOcT531Ty", nil)

				return userRepo
			},
			expectedResult: entity.LoginResponse{
				Success: true,
				Code:    http.StatusOK,
				Message: "OK",
				Token:   "publictoken",
				User: entity.User{
					Username: "foo",
					Name:     "Foo Bar",
				},
			},
			expectError: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			authController := controllers.NewAuth(controllers.AuthControllerParam{
				UserRP: tc.mockUserRepo(),
				AuthRP: tc.mockAuthRepo(),
			})

			result, err := authController.Login(tc.username, tc.password)

			assert.Equal(t, tc.expectedResult, result)
			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthenticateRequest(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tcs := []struct {
		name           string
		token          string
		mockAuthRepo   func() repository.AuthRepo
		mockUserRepo   func() repository.UserRepo
		expectedResult entity.User
		expectedError  error
	}{
		// Empty token
		{
			name:  "Empty token",
			token: "",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				return userRepo
			},
			expectedResult: entity.User{},
			expectedError:  errors.New("empty token"),
		},
		// Error parse token from repository
		{
			name:  "Error parse token from repository",
			token: "tokenfromheader",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				authRepo.EXPECT().ParseToken("tokenfromheader").Return(entity.TokenClaims{}, errors.New("some error"))

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				return userRepo
			},
			expectedResult: entity.User{},
			expectedError:  errors.New("some error"),
		},
		// User not found in our DB
		{
			name:  "User not found in our DB",
			token: "tokenfromheader",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				authRepo.EXPECT().ParseToken("tokenfromheader").Return(entity.TokenClaims{
					Name:     "Foo Bar",
					Username: "foobar",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				userRepo.EXPECT().GetByUsername("foobar").Return(entity.User{}, sql.ErrNoRows)

				return userRepo
			},
			expectedResult: entity.User{},
			expectedError:  sql.ErrNoRows,
		},
		// Successfully get user from token
		{
			name:  "Successfully get user from token",
			token: "tokenfromheader",
			mockAuthRepo: func() repository.AuthRepo {
				authRepo := repository.NewMockAuthRepo(mockController)

				authRepo.EXPECT().ParseToken("tokenfromheader").Return(entity.TokenClaims{
					Name:     "Foo Bar",
					Username: "foobar",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)

				return authRepo
			},
			mockUserRepo: func() repository.UserRepo {
				userRepo := repository.NewMockUserRepo(mockController)

				userRepo.EXPECT().GetByUsername("foobar").Return(entity.User{
					ID:        100,
					Username:  "foobar",
					Name:      "Foo Bar",
					Email:     "",
					Role:      "Super Admin",
					Status:    1,
					CreatedBy: "Test",
					UpdatedBy: "Test",
				}, nil)

				return userRepo
			},
			expectedResult: entity.User{
				ID:        100,
				Username:  "foobar",
				Name:      "Foo Bar",
				Email:     "",
				Role:      "Super Admin",
				Status:    1,
				CreatedBy: "Test",
				UpdatedBy: "Test",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			authController := controllers.NewAuth(controllers.AuthControllerParam{
				UserRP: tc.mockUserRepo(),
				AuthRP: tc.mockAuthRepo(),
			})

			result, err := authController.GetUserFromToken(tc.token)

			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
