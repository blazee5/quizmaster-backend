package domain

type Answer struct {
	Text       string `json:"text" validate:"required"`
	QuestionId int
	IsCorrect  bool `json:"is_correct" validate:"required"`
}

type UserAnswer struct {
	Id int `json:"id" validate:"required"`
}
