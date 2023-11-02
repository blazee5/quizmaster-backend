package models

type User struct {
	Id       int    `json:"id" db:"id" redis:"id"`
	Fio      string `json:"fio" db:"fio" redis:"fio"`
	Email    string `json:"email" db:"email" redis:"email"`
	Password string `json:"password" db:"password" redis:"password"`
}
