package domain

type Question struct {
	Title   string `json:"title" validate:"required"`
	Image   string
	Answers []Answer `json:"answers" validate:"required"`
	QuizId  int
}
