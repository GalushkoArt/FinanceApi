package model

import (
	"errors"
	"time"
)

type RefreshToken struct {
	UserId    string
	Token     string
	ExpiresAt time.Time
}

var (
	TokenNotFound = errors.New("token not found")
	TokenExpired  = errors.New("token expired")
)
