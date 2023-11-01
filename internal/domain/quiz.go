package domain

type Quiz struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description" validate:"required"`
	Image       string `form:"image"`
	UserId      int
	Questions   []Question `form:"questions" validate:"required"`
}
