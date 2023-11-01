package domain

type Question struct {
	Title   string   `json:"title" validate:"required"`
	Image   string   `json:"image" validate:"required"`
	Answers []Answer `json:"answers" validate:"required"`
	QuizId  int
}
