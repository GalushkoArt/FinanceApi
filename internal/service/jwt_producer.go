package service

import (
	"context"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"time"
)

const issuer = "financeapi.io"

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JwtProducer struct {
	expiryTimeout time.Duration
	hmacSecret    []byte
}

func NewJwtProducer(hmacSecret string, expiryTimeout time.Duration) *JwtProducer {
	return &JwtProducer{hmacSecret: []byte(hmacSecret), expiryTimeout: expiryTimeout}
}

func (s *JwtProducer) GetToken(ctx context.Context, id string, role model.Role) (string, error) {
	claims := &Claims{
		Role: string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    issuer,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiryTimeout)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.hmacSecret)
	if err != nil {
		utils.LogRequest(ctx, log.Error()).Str("from", "jwtProducer").Err(err).Msg("Failed to generate token string")
	}
	return tokenString, err
}
