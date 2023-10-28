package models

type User struct {
	Id       int    `json:"id" db:"id"`
	Fio      string `json:"fio" db:"fio"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
