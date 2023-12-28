package email

type Service interface {
	SendEmail(msg string) error
}
