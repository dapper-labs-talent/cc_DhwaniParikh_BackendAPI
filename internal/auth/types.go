package auth

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

type JwtClaims struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Id       uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}
