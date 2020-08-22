package repository

import (
	"github.com/nikolausreza131192/pos/entity"
	"log"
)

type AuthRepo interface {
	GenerateJWTToken(claims entity.TokenClaims) (string, error)
}

type authRepo struct {
	secretToken string
	jwtTokenLib entity.JWTToken
}

type AuthRepoParam struct {
	SecretToken string
	JWTTokenLib entity.JWTToken
}

func NewAuth(param AuthRepoParam) AuthRepo {
	return &authRepo{
		secretToken: param.SecretToken,
		jwtTokenLib: param.JWTTokenLib,
	}
}

func (r *authRepo) GenerateJWTToken(claims entity.TokenClaims) (string, error) {
	r.jwtTokenLib.SetClaims(claims)

	tokenString, err := r.jwtTokenLib.SignedString([]byte(r.secretToken))
	if err != nil {
		log.Println("GenerateJWTToken Error generate token", err)
		return "", err
	}

	return tokenString, nil
}
