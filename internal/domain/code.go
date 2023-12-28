package domain

type VerificationCode struct {
	Email string `json:"email" validate:"required"`
	Type  string `json:"type" validate:"required,oneof=email password"`
}
