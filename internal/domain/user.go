package domain

type User struct {
	Id     int    `json:"id"`
	Fio    string `json:"fio"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

type UpdateUser struct {
	Fio string `json:"fio" validate:"required,min=4"`
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
