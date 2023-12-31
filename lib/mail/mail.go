package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

const (
	EmailConfirmationType = "confirm"
	ResetEmailType        = "email"
	ResetPasswordType     = "password"
)

func SendMail(emailType, username, email, code string) error {
	templates := map[string]string{
		EmailConfirmationType: "../lib/templates/email-confirm.html",
		ResetEmailType:        "../lib/templates/reset-email.html",
		ResetPasswordType:     "../lib/templates/reset-password.html",
	}

	t, err := template.ParseFiles(templates[emailType])
	if err != nil {
		return err
	}

	var body bytes.Buffer

	var link string

	switch emailType {
	case ResetPasswordType:
		link = fmt.Sprintf("https://quizer-opal.vercel.app/user/reset/password/%s", code)
	case ResetEmailType:
		link = fmt.Sprintf("https://quizer-opal.vercel.app/user/reset/email/%s", code)
	}

	if err := t.Execute(&body, map[string]any{"Link": link, "Username": username}); err != nil {
		return err
	}

	message := []byte("Subject: Account Activation\r\n" +
		"From: " + os.Getenv("SMTP_FROM") + "\r\n" +
		"To: " + email + "\r\n" +
		"MIME-version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body.String())

	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	err = smtp.SendMail(os.Getenv("SMTP_ADDR"), auth, os.Getenv("SMTP_FROM"), []string{email}, message)
	if err != nil {
		return err
	}

	return nil
}
