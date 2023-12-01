package domain

type SubmitResult struct {
	AttemptID int `json:"attempt_id" validate:"required"`
}
