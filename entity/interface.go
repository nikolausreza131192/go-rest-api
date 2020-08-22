package entity

type (
	JWTToken interface {
		SetClaims(claims TokenClaims)
		SignedString(token []byte) (string, error)
	}
)
