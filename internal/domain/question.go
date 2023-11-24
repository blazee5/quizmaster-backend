package domain

type Question struct {
	Title  string `json:"title"`
	Type   string `json:"type" validate:"required,oneof=choice input"`
	QuizID int
}

type ChangeQuestionOrder struct {
	FirstOrderID  float64 `json:"first_order_id" validate:"required"`
	SecondOrderID float64 `json:"second_order_id" validate:"required"`
	QuestionID    int     `json:"question_id" validate:"required"`
}
