package models

import "time"

type Result struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	QuizID      int       `json:"quiz_id" db:"quiz_id"`
	Score       int       `json:"score" db:"score"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UsersResult struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username,omitempty" db:"username"`
	Score     int       `json:"score" db:"score"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserResult struct {
	Quiz        Quiz      `json:"quiz" db:"quiz"`
	Score       int       `json:"score" db:"score"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
