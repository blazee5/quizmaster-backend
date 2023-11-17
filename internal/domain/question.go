package domain

type Question struct {
	Title  string `json:"title"`
	Type   string `json:"type" validate:"required,oneof=choice input"`
	QuizId int
}
