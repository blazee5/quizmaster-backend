package domain

type Question struct {
	Title  string `json:"title" validate:"required"`
	Type   string `json:"type" validate:"required,oneof=choice input"`
	QuizId int
}
