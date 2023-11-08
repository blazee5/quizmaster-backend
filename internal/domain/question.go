package domain

type Question struct {
	Title   string   `json:"title" validate:"required"`
	Answers []Answer `json:"answers" validate:"required"`
	Image   string   `json:"image"`
	Type    string   `json:"type" validate:"required,oneof=choice,input"`
	QuizId  int
}
