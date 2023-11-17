package domain

type Answer struct {
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	QuestionId int
}

type UserAnswer struct {
	Id int `json:"id" validate:"required"`
}
