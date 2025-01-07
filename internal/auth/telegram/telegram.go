package telegram

type Telegram struct {
	token string
}

func New(token string) *Telegram {
	return &Telegram{
		token: token,
	}
}
