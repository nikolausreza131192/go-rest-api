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
