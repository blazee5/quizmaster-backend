package domain

type Quiz struct {
	Title       string     `form:"title" validate:"required"`
	Description string     `form:"description" validate:"required"`
	Questions   []Question `validate:"required"`
	Image       string
	UserId      int
}
