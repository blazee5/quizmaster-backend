package domain

type Email struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Message string `json:"message"`
}
