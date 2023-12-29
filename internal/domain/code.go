package domain

type VerificationCode struct {
	Email string `json:"email"`
	Type  string `json:"type" validate:"required,oneof=email password"`
}
