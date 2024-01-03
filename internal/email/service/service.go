package service

import (
	"encoding/json"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/mail"
	"go.uber.org/zap"
)

type Service struct {
	log *zap.SugaredLogger
}

func NewService(log *zap.SugaredLogger) *Service {
	return &Service{log: log}
}

func (s *Service) SendEmail(msg string) error {
	var email domain.Email

	err := json.Unmarshal([]byte(msg), &email)
	if err != nil {
		s.log.Infof("error while unmarshal message: %v", err)
		return err
	}

	err = mail.SendMail(email.Type, email.Username, email.To, email.Code)

	if err != nil {
		s.log.Infof("error while send email: %v", err)
		return err
	}

	return nil
}
