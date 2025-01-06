package jwt

import "time"

type Auth struct {
	key      []byte
	tokenTTL time.Duration
}

func NewAuth(secretKey string) *Auth {
	return &Auth{
		key: []byte(secretKey),
	}
}
