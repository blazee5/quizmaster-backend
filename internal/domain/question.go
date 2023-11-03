package domain

type Question struct {
	Title   string   `json:"title" validate:"required"`
	Answers []Answer `json:"answers" validate:"required"`
	Image   string
	QuizId  int
}
