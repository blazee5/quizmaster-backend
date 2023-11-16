package models

type Answer struct {
	Id         int     `json:"id" db:"id"`
	Text       string  `json:"text" db:"text"`
	QuestionId int     `json:"question_id" db:"question_id"`
	IsCorrect  bool    `json:"is_correct,omitempty" db:"is_correct"`
	OrderId    float64 `json:"order_id" db:"order_id"`
}
