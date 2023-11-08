package domain

type Result struct {
	Answers map[int]interface{} `json:"answers" validate:"required"`
}
