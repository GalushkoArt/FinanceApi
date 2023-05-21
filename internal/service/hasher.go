package service

import (
	"crypto/sha256"
	"fmt"
	"github.com/rs/zerolog/log"
)

type Hasher struct {
	salt []byte
}

func NewHasher(salt string) *Hasher {
	if len(salt) == 0 {
		log.Panic().Msg("salt for hasher cannot be empty!")
	}
	return &Hasher{salt: []byte(salt)}
}

func (h *Hasher) Hash(input string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(input)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(h.salt)), nil
}
