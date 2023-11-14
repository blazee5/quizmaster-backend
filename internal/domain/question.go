package domain

type Question struct {
	Title  string `form:"title" validate:"required"`
	Image  string `form:"image"`
	Type   string `form:"type" validate:"required,oneof=choice input"`
	QuizId int
}
