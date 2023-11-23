package domain

type Answer struct {
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	QuestionID int
}

type UserAnswer struct {
	ID int `json:"id" validate:"required"`
}
