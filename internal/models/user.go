package models

type User struct {
	ID         int    `json:"id" db:"id" redis:"id"`
	Username   string `json:"username" db:"username" redis:"username"`
	Email      string `json:"email" db:"email" redis:"email"`
	Password   string `json:"password" db:"password" redis:"password"`
	Avatar     string `json:"avatar" db:"avatar" redis:"avatar"`
	RoleID     int    `json:"role_id" db:"role_id" redis:"role_id"`
	IsVerified bool   `json:"is_verified" db:"is_verified" redis:"is_verified"`
}

type ShortUser struct {
	ID       int    `json:"id" db:"id" redis:"id"`
	Username string `json:"username" db:"username" redis:"username"`
	Email    string `json:"email" db:"email" redis:"email"`
	Avatar   string `json:"avatar" db:"avatar" redis:"avatar"`
}

type UserInfo struct {
	User    ShortUser    `json:"user"`
	Quizzes []Quiz       `json:"quizzes"`
	Results []UserResult `json:"results"`
}
