package models

type TelegramCallback struct {
	ID        int64   `form:"id" binding:"required"`
	FirstName string  `form:"first_name" binding:"required"`
	Username  string  `form:"username" binding:"required"`
	PhotoURL  *string `form:"photo_url"`
	AuthDate  int64   `form:"auth_date"`
	Hash      string  `form:"hash" binding:"required"`
	Next      string  `form:"next"`
}
