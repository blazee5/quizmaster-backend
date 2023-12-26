package domain

type Question struct {
	Title   string `json:"title"`
	Type    string `json:"type" validate:"required,oneof=choice input"`
	OrderID int    `json:"order_id"`
	QuizID  int
}

type QuestionOrder struct {
	QuestionID int `json:"question_id" validate:"required"`
	OrderID    int `json:"order_id" validate:"required"`
}
