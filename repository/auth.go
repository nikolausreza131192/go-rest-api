package repository

import (
	"errors"
	"fmt"
	"github.com/nikolausreza131192/pos/entity"
	"log"
	"sync"
)

type AuthRepo interface {
	InitPermissionMap()
	GenerateJWTToken(claims entity.TokenClaims) (string, error)
	ParseToken(token string) (entity.TokenClaims, error)
	GetPermissionsByRole(roleName string) []string
}

type authRepo struct {
	secretToken   string
	jwtTokenLib   entity.JWTToken
	permissionMap *sync.Map
}

type AuthRepoParam struct {
	SecretToken string
	JWTTokenLib entity.JWTToken
}

func NewAuth(param AuthRepoParam) AuthRepo {
	r := &authRepo{
		secretToken: param.SecretToken,
		jwtTokenLib: param.JWTTokenLib,
	}
	r.InitPermissionMap()

	return r
}

func (r *authRepo) InitPermissionMap() {
	// TODO : use DB instead of hardcoded map
	permissionResources := map[string][]string{
		"get list items": {
			"Super Admin",
		},
		"get item details": {
			"Super Admin",
			"Admin",
		},
	}

	r.permissionMap = &sync.Map{}
	for permissionName, authorizedRoles := range permissionResources {
		r.permissionMap.Store(fmt.Sprintf("%s", permissionName), authorizedRoles)
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

func (r *authRepo) ParseToken(token string) (entity.TokenClaims, error) {
	var tokenClaims entity.TokenClaims
	if token == "" {
		log.Println("ParseToken empty token param")
		return tokenClaims, errors.New("empty token")
	}

	tokenClaims, err := r.jwtTokenLib.ParseToken(token, r.secretToken)
	if err != nil {
		log.Println("ParseToken error parse token", err)
		return tokenClaims, err
	}

	return tokenClaims, nil
}

func (r *authRepo) GetPermissionsByRole(roleName string) []string {
	var permissions []string

	r.permissionMap.Range(func(key, value interface{}) bool {
		roles, ok := value.([]string)
		if !ok {
			log.Println("GetPermissionsByRole incorrect format in permission map", roleName)

			// Continue to next loop
			return true
		}
		permission, ok := key.(string)
		if !ok {
			log.Println("GetPermissionsByRole permission is not string")

			// Continue to next loop
			return true
		}

		for _, role := range roles {
			if role == roleName {
				permissions = append(permissions, permission)
			}
		}

		return true
	})

	return permissions
}
