package entity

type Role struct {
	ID          int64  `json:"id"`
	Level       int    `json:"level"`
	Name        string `json:"title"`
	Description string `json:"description"`
}
