package domain

type Email struct {
	Type     string `json:"type"`
	To       string `json:"to"`
	Username string `json:"username"`
	Code     string `json:"code"`
}
