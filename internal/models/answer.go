package models

type Answer struct {
	ID         int    `json:"id" db:"id"`
	Text       string `json:"text" db:"text"`
	QuestionID int    `json:"question_id" db:"question_id"`
	IsCorrect  bool   `json:"is_correct,omitempty" db:"is_correct"`
	OrderID    int    `json:"order_id" db:"order_id"`
}

type AnswerInfo struct {
	ID         int    `json:"id" db:"id"`
	Text       string `json:"text" db:"text"`
	QuestionID int    `json:"question_id" db:"question_id"`
	IsCorrect  bool   `json:"is_correct" db:"is_correct"`
	OrderID    int    `json:"order_id" db:"order_id"`
}
