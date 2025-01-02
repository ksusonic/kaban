package models

type User struct {
	ID         int
	Username   string
	TelegramID int64
	FirstName  string
	LastName   string
	AvatarURL  *string
}
