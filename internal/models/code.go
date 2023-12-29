package models

import "time"

type VerificationCode struct {
	ID         int       `json:"id" db:"id"`
	Type       string    `json:"type" db:"type"`
	Code       string    `json:"code" db:"code"`
	UserID     int       `json:"user_id" db:"user_id"`
	ExpireDate time.Time `json:"expire_date" db:"expire_date"`
}
