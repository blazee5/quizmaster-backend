package domain

type Result struct {
	QuizId  int           `json:"quiz_id"`
	Answers map[int][]int `json:"answers"`
}
