package domain

type VerificationCode struct {
	Email string `json:"email" validate:"required,email"`
}
