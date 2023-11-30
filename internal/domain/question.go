package domain

type Question struct {
	Title   string `json:"title"`
	Type    string `json:"type" validate:"required,oneof=choice input"`
	OrderID int    `json:"order_id"`
	QuizID  int
}

type OrderQuestionItem struct {
	QuestionID int `json:"question_id" validate:"required"`
	OrderID    int `json:"order_id" validate:"required"`
}

type ChangeQuestionOrder struct {
	Orders []OrderQuestionItem `json:"orders" validate:"required"`
}
