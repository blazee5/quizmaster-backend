package models

type User struct {
	Id       int    `json:"id" db:"id" redis:"id"`
	Username string `json:"username" db:"username" redis:"username"`
	Email    string `json:"email" db:"email" redis:"email"`
	Password string `json:"password" db:"password" redis:"password"`
	Avatar   string `json:"avatar" db:"avatar" redis:"avatar"`
}
