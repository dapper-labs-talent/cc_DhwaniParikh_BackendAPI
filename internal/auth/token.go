package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"server/internal/config"
	"server/internal/models"
	"time"

	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

type TokenService struct {
	config config.AuthConfig
	logger *zerolog.Logger
}

var JWT_SECRET string

func NewTokenService(c config.AuthConfig, l *zerolog.Logger) *TokenService {
	return &TokenService{
		config: c,
		logger: l,
	}
}

func (ts *TokenService) GenerateAPITokens(id uuid.UUID, email string) (tokens *models.TokenResponse, err error) {
	result := &models.TokenResponse{}
	// Generate a new access token
	result.Token, err = generateToken(id, email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func generateToken(id uuid.UUID, email string) (string, error) {

	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	key := []byte(JWT_SECRET)
	expirationTime := time.Now().Add(7 * 24 * 60 * time.Minute)
	claims := &JwtClaims{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	UnsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := UnsignedToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return SignedToken, nil
}
