package auth

import (
	"time"

	"github.com/ksusonic/kanban/internal/auth/telegram"
)

type Auth struct {
	*telegram.Telegram
	key      []byte
	tokenTTL time.Duration
}

func NewAuth(
	secretKey string,
	tgToken string,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		Telegram: telegram.New(tgToken),
		key:      []byte(secretKey),
		tokenTTL: tokenTTL,
	}
}
