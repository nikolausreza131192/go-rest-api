package repository_test

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateJWTToken(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tcs := []struct {
		name             string
		param            entity.TokenClaims
		mockJWT          func() *repository.MockJWTToken
		expectedResponse string
		expectedError    error
	}{
		{
			name:  "Error signed string",
			param: entity.TokenClaims{},
			mockJWT: func() *repository.MockJWTToken {
				jwtToken := repository.NewMockJWTToken(mockController)

				jwtToken.EXPECT().SetClaims(gomock.Any())
				jwtToken.EXPECT().SignedString([]byte("secrettoken")).Return("", errors.New("test error"))

				return jwtToken
			},
			expectedResponse: "",
			expectedError:    errors.New("test error"),
		},
		{
			name: "Success generate token",
			param: entity.TokenClaims{
				Username: "foobar",
				Name:     "Foo Bar",
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Unix(),
				},
			},
			mockJWT: func() *repository.MockJWTToken {
				jwtToken := repository.NewMockJWTToken(mockController)

				jwtToken.EXPECT().SetClaims(gomock.Any())
				jwtToken.EXPECT().SignedString([]byte("secrettoken")).Return("generatedtoken", nil)

				return jwtToken
			},
			expectedResponse: "generatedtoken",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			param := repository.AuthRepoParam{
				SecretToken: "secrettoken",
				JWTTokenLib: tc.mockJWT(),
			}

			repo := repository.NewAuth(param)
			token, err := repo.GenerateJWTToken(tc.param)

			assert.Equal(t, tc.expectedResponse, token)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestParseToken(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	tcs := []struct {
		name             string
		param            string
		mockJWT          func() *repository.MockJWTToken
		expectedResponse entity.TokenClaims
		expectedError    error
	}{
		// Token is empty
		{
			name:  "Token is empty",
			param: "",
			mockJWT: func() *repository.MockJWTToken {
				jwtToken := repository.NewMockJWTToken(mockController)

				return jwtToken
			},
			expectedResponse: entity.TokenClaims{},
			expectedError:    errors.New("empty token"),
		},
		// Error parse token
		{
			name:  "Error parse token",
			param: "tokenfromheader",
			mockJWT: func() *repository.MockJWTToken {
				jwtToken := repository.NewMockJWTToken(mockController)

				jwtToken.EXPECT().ParseToken("tokenfromheader", "secrettoken").Return(entity.TokenClaims{}, errors.New("some error"))

				return jwtToken
			},
			expectedResponse: entity.TokenClaims{},
			expectedError:    errors.New("some error"),
		},
		// Success
		{
			name:  "Success",
			param: "tokenfromheader",
			mockJWT: func() *repository.MockJWTToken {
				jwtToken := repository.NewMockJWTToken(mockController)

				jwtToken.EXPECT().ParseToken("tokenfromheader", "secrettoken").Return(entity.TokenClaims{
					Username: "foobar",
					Name:     "Foo Bar",
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC).Unix(),
					},
				}, nil)

				return jwtToken
			},
			expectedResponse: entity.TokenClaims{
				Username: "foobar",
				Name:     "Foo Bar",
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC).Unix(),
				},
			},
			expectedError: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			param := repository.AuthRepoParam{
				SecretToken: "secrettoken",
				JWTTokenLib: tc.mockJWT(),
			}

			repo := repository.NewAuth(param)
			tokenClaims, err := repo.ParseToken(tc.param)

			assert.Equal(t, tc.expectedResponse, tokenClaims)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetPermissionsByRole(t *testing.T) {
	tcs := []struct {
		name            string
		param           string
		expectedResults []string
	}{
		{
			name:  "Success get permissions for super admin",
			param: "Super Admin",
			expectedResults: []string{
				"get list items",
				"get item details",
			},
		},
		{
			name:  "Success get permissions for admin",
			param: "Admin",
			expectedResults: []string{
				"get list items",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := repository.NewAuth(repository.AuthRepoParam{})

			permissions := authRepo.GetPermissionsByRole(tc.param)

			assert.Equal(t, tc.expectedResults, permissions)
		})
	}
}
