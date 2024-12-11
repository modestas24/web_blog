package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	RoleID    int64     `json:"role_id"`
	Role      Role      `json:"-"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  Password  `json:"-"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Password struct {
	Raw  string
	Hash []byte
}

func (password *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	password.Raw = text
	password.Hash = hash

	return nil
}

func (password *Password) Compare(raw []byte) error {
	return bcrypt.CompareHashAndPassword(password.Hash, raw)
}
