package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/answer"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/question"
	"github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/internal/result"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.uber.org/zap"
	"strings"
)

type Service struct {
	log          *zap.SugaredLogger
	repo         result.Repository
	quizRepo     quiz.Repository
	questionRepo question.Repository
	answerRepo   answer.Repository
}

func NewService(log *zap.SugaredLogger, repo result.Repository, quizRepo quiz.Repository, questionRepo question.Repository, answerRepo answer.Repository) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo, questionRepo: questionRepo, answerRepo: answerRepo}
}

func (s *Service) NewResult(ctx context.Context, userID int, quizID int) (int, error) {
	_, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return 0, err
	}

	return s.repo.NewResult(ctx, userID, quizID)
}

func (s *Service) SaveUserAnswer(ctx context.Context, userID, quizID int, input domain.UserAnswer) error {
	if _, err := s.quizRepo.GetByID(ctx, quizID); err != nil {
		return err
	}

	attempt, err := s.repo.GetByID(ctx, input.AttemptID)

	if err != nil {
		return err
	}

	if attempt.UserID != userID || attempt.IsCompleted {
		return http_errors.PermissionDenied
	}

	question, err := s.questionRepo.GetQuestionByID(ctx, input.QuestionID)

	if err != nil {
		return err
	}

	if question.QuizID != quizID {
		return http_errors.WrongArgument
	}

	answerExists, err := s.repo.GetUserAnswerByID(ctx, input.AnswerID, input.AttemptID)

	if err != nil {
		return err
	}

	if answerExists {
		return http_errors.PermissionDenied
	}

	if err := s.repo.SaveUserAnswer(ctx, userID, input.QuestionID, input.AnswerID, input.AttemptID, input.AnswerText); err != nil {
		return err
	}

	if question.Type == "choice" {
		err := s.ProcessChoiceAnswer(ctx, question.ID, userID, input)

		if err != nil {
			return err
		}
	} else {
		err := s.ProcessInputAnswer(ctx, question.ID, userID, input)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ProcessChoiceAnswer(ctx context.Context, questionID, userID int, input domain.UserAnswer) error {
	answer, err := s.answerRepo.GetByID(ctx, input.AnswerID)

	if err != nil {
		return err
	}

	if answer.QuestionID != questionID {
		return http_errors.WrongArgument
	}

	if answer.IsCorrect {
		err := s.repo.UpdateResult(ctx, input.AttemptID, userID, 1)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) ProcessInputAnswer(ctx context.Context, questionID, userID int, input domain.UserAnswer) error {
	correctAnswers, err := s.answerRepo.GetAnswersByQuestionID(ctx, questionID)

	if err != nil {
		return err
	}

	for _, ans := range correctAnswers {
		if strings.ToLower(ans.Text) == strings.ToLower(input.AnswerText) {
			err := s.repo.UpdateResult(ctx, input.AttemptID, userID, 1)

			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}

func (s *Service) GetResultsByQuizID(ctx context.Context, quizID int) ([]models.UsersResult, error) {
	return s.repo.GetByQuizID(ctx, quizID)
}

func (s *Service) SubmitResult(ctx context.Context, userID, quizID int, input domain.SubmitResult) (models.UsersResult, error) {
	if _, err := s.quizRepo.GetByID(ctx, quizID); err != nil {
		return models.UsersResult{}, err
	}

	return s.repo.SubmitResult(ctx, userID, input.AttemptID)
}
