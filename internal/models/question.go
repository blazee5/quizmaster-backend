package models

type Question struct {
	ID      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Image   string `json:"image" db:"image"`
	QuizID  int    `json:"quiz_id" db:"quiz_id"`
	Type    string `json:"type" db:"type"`
	OrderID int    `json:"order_id" db:"order_id"`
}

type QuestionWithAnswers struct {
	ID      int      `json:"id" db:"id"`
	Title   string   `json:"title" db:"title"`
	Image   string   `json:"image" db:"image"`
	QuizID  int      `json:"quiz_id" db:"quiz_id"`
	Type    string   `json:"type" db:"type"`
	OrderID int      `json:"order_id" db:"order_id"`
	Answers []Answer `json:"answers" db:"answers"`
}
