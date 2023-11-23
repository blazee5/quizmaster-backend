package models

import "time"

type Quiz struct {
	ID          int       `json:"id" db:"id" redis:"id"`
	Title       string    `json:"title" db:"title" redis:"title"`
	Description string    `json:"description" db:"description" redis:"description"`
	Image       string    `json:"image" db:"image" redis:"image"`
	UserID      int       `json:"user_id" db:"user_id" redis:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" redis:"created_at"`
}
