package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/nikolausreza131192/pos/config"
	"github.com/nikolausreza131192/pos/entity"
)

type token struct {
	signingMethod string
	jwtToken      *jwt.Token
}

func initJWTToken(conf config.JWTConfig) entity.JWTToken {
	return &token{
		signingMethod: conf.SigningMethod,
		jwtToken:      jwt.New(jwt.GetSigningMethod(conf.SigningMethod)),
	}
}

func (t *token) SetClaims(claims entity.TokenClaims) {
	t.jwtToken.Header = map[string]interface{}{
		"typ": "JWT",
		"alg": jwt.GetSigningMethod(t.signingMethod).Alg(),
	}
	t.jwtToken.Claims = claims
	t.jwtToken.Method = jwt.GetSigningMethod(t.signingMethod)
}

func (t *token) SignedString(method []byte) (string, error) {
	return t.jwtToken.SignedString(method)
}
