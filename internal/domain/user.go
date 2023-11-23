package domain

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UpdateUser struct {
	Username string `json:"username" validate:"required"`
}

type UserQuiz struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
	Image       string     `json:"image"`
}

type UserResult struct {
	Quizzes []Quiz `json:"quizzes"`
}
