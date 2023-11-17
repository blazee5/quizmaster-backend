package models

type Question struct {
	Id      int      `json:"id" db:"id"`
	Title   string   `json:"title" db:"title"`
	Image   string   `json:"image" db:"image"`
	Answers []Answer `json:"answers"`
	QuizId  int      `json:"quiz_id" db:"quiz_id"`
	Type    string   `json:"type" db:"type"`
	OrderId float64  `json:"order_id" db:"order_id"`
}

type QuestionWithAnswers struct {
	Id      int          `json:"id" db:"id"`
	Title   string       `json:"title" db:"title"`
	Image   string       `json:"image" db:"image"`
	Answers []AnswerInfo `json:"answers" db:"answers"`
	QuizId  int          `json:"quiz_id" db:"quiz_id"`
	Type    string       `json:"type" db:"type"`
	OrderId float64      `json:"order_id" db:"order_id"`
}
