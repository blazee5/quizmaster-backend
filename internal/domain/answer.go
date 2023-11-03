package domain

type Answer struct {
	Text       string `json:"text" validate:"required"`
	IsCorrect  bool   `json:"is_correct" validate:"required"`
	QuestionId int
}

type UserAnswer struct {
	Id int `json:"id" validate:"required"`
}
