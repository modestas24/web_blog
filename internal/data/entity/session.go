package entity

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiredAt time.Time `json:"expired_at"`
}
