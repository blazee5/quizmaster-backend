package models

import "time"

type Result struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	QuizId    int       `json:"quiz_id" db:"quiz_id"`
	Score     int       `json:"score" db:"score"`
	Percent   int       `json:"percent" db:"percent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserResult struct {
	Quiz      Quiz      `json:"quiz" db:"quiz"`
	Score     int       `json:"score" db:"score"`
	Percent   int       `json:"percent" db:"percent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
