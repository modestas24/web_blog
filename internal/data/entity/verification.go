package entity

import (
	"time"

	"github.com/google/uuid"
)

type Verification struct {
	UUID      uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	ExpiredAt time.Time `json:"expired_at"`
}
