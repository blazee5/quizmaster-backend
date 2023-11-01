package domain

type Quiz struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image"`
	UserId      int
	Questions   []Question `json:"questions" validate:"required"`
}
