package jwt

import "time"

const issuer = "Kkaban board"

type Auth struct {
	key      []byte
	tokenTTL time.Duration
}

func NewAuth(secretKey string) *Auth {
	return &Auth{
		key: []byte(secretKey),
	}
}
