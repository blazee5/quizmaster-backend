package models

type Question struct {
	Id      int      `json:"id" db:"id"`
	Title   string   `json:"title" db:"title"`
	Image   string   `json:"image" db:"image"`
	Answers []Answer `json:"answers" db:"answers"`
	QuizId  int      `json:"quiz_id" db:"quiz_id"`
}
