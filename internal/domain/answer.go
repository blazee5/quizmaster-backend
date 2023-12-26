package domain

type Answer struct {
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderID    int    `json:"order_id"`
	QuestionID int
}

type UserAnswer struct {
	AttemptID  int    `json:"attempt_id" validate:"required"`
	QuestionID int    `json:"question_id" validate:"required"`
	AnswerID   int    `json:"answer_id"`
	AnswerText string `json:"answer_text"`
}

type AnswerOrder struct {
	AnswerID int `json:"answer_id" validate:"required"`
	OrderID  int `json:"order_id" validate:"required"`
}
