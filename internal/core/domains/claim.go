package domains

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
