package models

type User struct {
	ID         int
	Username   string
	FirstName  string
	TelegramID *int64
	LastName   *string
	AvatarURL  *string
}
