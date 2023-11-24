package domain

type Answer struct {
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	QuestionID int
}

type UserAnswer struct {
	ID int `json:"id" validate:"required"`
}

type ChangeAnswerOrder struct {
	From     int `json:"from" validate:"required"`
	To       int `json:"to" validate:"required"`
	AnswerID int `json:"answer_id" validate:"required"`
}
