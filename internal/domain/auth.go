package domain

type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ResetEmailRequest struct {
	Code string `json:"code" validate:"required"`
}

type ResetPasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
