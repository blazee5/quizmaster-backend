package models

type Answer struct {
	ID         int    `json:"id" db:"id"`
	Text       string `json:"text" db:"text"`
	QuestionID int    `json:"question_id" db:"question_id"`
	IsCorrect  bool   `json:"is_correct" db:"is_correct"`
	OrderID    int    `json:"order_id" db:"order_id"`
}

type AnswerInfo struct {
	ID         int    `json:"id" db:"id"`
	Text       string `json:"text" db:"text"`
	QuestionID int    `json:"question_id" db:"question_id"`
	OrderID    int    `json:"order_id" db:"order_id"`
}

type UserAnswer struct {
	ID         int    `json:"id" db:"id"`
	UserID     int    `json:"user_id" db:"user_id"`
	QuestionID int    `json:"question_id" db:"question_id"`
	AnswerID   int    `json:"answer_id" db:"answer_id"`
	ResultID   int    `json:"result_id" db:"result_id"`
	Text       string `json:"text" db:"text"`
}
