package service

import (
	"FinanceApi/internal/model"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type JwtParser struct {
	hmacSecret []byte
}

func NewJwtParser(hmacSecret string) *JwtParser {
	return &JwtParser{hmacSecret: []byte(hmacSecret)}
}

func (s *JwtParser) ParseToken(token string) (string, model.Role, error) {
	claims := Claims{}
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return "", model.ClientRole, err
	}
	if !t.Valid || claims.Issuer != issuer {
		return "", model.ClientRole, errors.New("invalid token")
	}
	return claims.Subject, model.Role(claims.Role), nil
}
