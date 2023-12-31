package service

import (
	"context"
	answerRepo "github.com/blazee5/quizmaster-backend/internal/answer"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	questionRepo "github.com/blazee5/quizmaster-backend/internal/question"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	log          *zap.SugaredLogger
	repo         answerRepo.Repository
	quizRepo     quizRepo.Repository
	questionRepo questionRepo.Repository
	tracer       trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo answerRepo.Repository, quizRepo quizRepo.Repository, questionRepo questionRepo.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo, questionRepo: questionRepo, tracer: tracer}
}

func (s *Service) Create(ctx context.Context, userID, quizID, questionID int) (int, error) {
	ctx, span := s.tracer.Start(ctx, "answerService.Create")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return 0, err
	}

	question, err := s.questionRepo.GetQuestionByID(ctx, questionID)

	if err != nil {
		return 0, err
	}

	if quiz.UserID != userID || question.QuizID != quizID {
		return 0, http_errors.ErrPermissionDenied
	}

	return s.repo.Create(ctx, questionID)
}

func (s *Service) GetByQuestionID(ctx context.Context, quizID, questionID int) ([]models.AnswerInfo, error) {
	ctx, span := s.tracer.Start(ctx, "answerService.GetByQuestionID")
	defer span.End()

	_, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return nil, err
	}

	question, err := s.questionRepo.GetQuestionByID(ctx, questionID)

	if err != nil {
		return nil, err
	}

	if question.Type == "input" {
		return nil, http_errors.ErrPermissionDenied
	}

	return s.repo.GetAnswersInfoByQuestionID(ctx, questionID)
}

func (s *Service) Update(ctx context.Context, answerID, userID, quizID, questionID int, input domain.Answer) error {
	ctx, span := s.tracer.Start(ctx, "answerService.Update")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, questionID, answerID)

	if err != nil {
		return err
	}

	return s.repo.Update(ctx, answerID, input)
}

func (s *Service) Delete(ctx context.Context, answerID, userID, quizID, questionID int) error {
	ctx, span := s.tracer.Start(ctx, "answerService.Delete")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, questionID, answerID)

	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, answerID)
}

func (s *Service) ChangeOrder(ctx context.Context, userID, quizID, questionID int, input domain.AnswerOrder) error {
	ctx, span := s.tracer.Start(ctx, "answerService.ChangeOrder")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, questionID, input.AnswerID)

	if err != nil {
		return err
	}

	err = s.repo.ChangeOrder(ctx, questionID, input)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) checkPermissions(ctx context.Context, userID, quizID, questionID, answerID int) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	question, err := s.questionRepo.GetQuestionByID(ctx, questionID)

	if err != nil {
		return err
	}

	answer, err := s.repo.GetByID(ctx, answerID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID || question.QuizID != quizID || answer.QuestionID != questionID {
		return http_errors.ErrPermissionDenied
	}

	return nil
}
