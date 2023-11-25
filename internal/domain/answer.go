package domain

type CreateAnswer struct {
	OrderID int `json:"order_id" validate:"required"`
}

type Answer struct {
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderID    int    `json:"order_id" validate:"required"`
	QuestionID int
}

type UserAnswer struct {
	ID int `json:"id" validate:"required"`
}

type OrderAnswerItem struct {
	AnswerID int `json:"answer_id" validate:"required"`
	OrderID  int `json:"order_id" validate:"required"`
}

type ChangeAnswerOrder struct {
	Orders []OrderAnswerItem `json:"orders" validate:"required"`
}
