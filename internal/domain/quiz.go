package domain

type Quiz struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description" validate:"required"`
	Image       string
	UserId      int
	Questions   []Question
}
