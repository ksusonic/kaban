package models

import "time"

type JWTToken struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
