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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strings"
)

type Service struct {
	log          *zap.SugaredLogger
	repo         result.Repository
	quizRepo     quiz.Repository
	questionRepo question.Repository
	answerRepo   answer.Repository
	tracer       trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo result.Repository, quizRepo quiz.Repository, questionRepo question.Repository, answerRepo answer.Repository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo, questionRepo: questionRepo, answerRepo: answerRepo, tracer: tracer}
}

func (s *Service) NewResult(ctx context.Context, userID int, quizID int) (int, error) {
	ctx, span := s.tracer.Start(ctx, "resultService.NewResult")
	defer span.End()

	_, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return 0, err
	}

	return s.repo.NewResult(ctx, userID, quizID)
}

func (s *Service) SaveUserAnswer(ctx context.Context, userID, quizID int, input domain.UserAnswer) error {
	ctx, span := s.tracer.Start(ctx, "resultService.SaveUserAnswer")
	defer span.End()

	if _, err := s.quizRepo.GetByID(ctx, quizID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	attempt, err := s.repo.GetByID(ctx, input.AttemptID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if attempt.UserID != userID || attempt.IsCompleted {
		return http_errors.ErrPermissionDenied
	}

	question, err := s.questionRepo.GetQuestionByID(ctx, input.QuestionID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if question.QuizID != quizID {
		return http_errors.ErrWrongArgument
	}

	answerExists, err := s.repo.GetUserAnswerByID(ctx, input.AnswerID, input.AttemptID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if answerExists {
		return http_errors.ErrPermissionDenied
	}

	if err := s.repo.SaveUserAnswer(ctx, userID, input.QuestionID, input.AnswerID, input.AttemptID, input.AnswerText); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if question.Type == "choice" {
		if input.AnswerID == 0 {
			return http_errors.ErrWrongArgument
		}

		err := s.ProcessChoiceAnswer(ctx, question.ID, userID, input)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return err
		}
	} else {
		if input.AnswerText == "" {
			return http_errors.ErrWrongArgument
		}

		err := s.ProcessInputAnswer(ctx, question.ID, userID, input)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return err
		}
	}

	return nil
}

func (s *Service) ProcessChoiceAnswer(ctx context.Context, questionID, userID int, input domain.UserAnswer) error {
	ctx, span := s.tracer.Start(ctx, "resultService.ProcessChoiceAnswer")
	defer span.End()

	answer, err := s.answerRepo.GetByID(ctx, input.AnswerID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if answer.QuestionID != questionID {
		return http_errors.ErrWrongArgument
	}

	if answer.IsCorrect {
		err := s.repo.UpdateResult(ctx, input.AttemptID, userID, 1)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return err
		}
	}

	return nil
}

func (s *Service) ProcessInputAnswer(ctx context.Context, questionID, userID int, input domain.UserAnswer) error {
	ctx, span := s.tracer.Start(ctx, "resultService.ProcessInputAnswer")
	defer span.End()

	correctAnswers, err := s.answerRepo.GetAnswersByQuestionID(ctx, questionID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	for _, ans := range correctAnswers {
		if strings.ToLower(ans.Text) == strings.ToLower(input.AnswerText) {
			err := s.repo.UpdateResult(ctx, input.AttemptID, userID, 1)

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				return err
			}

			break
		}
	}

	return nil
}

func (s *Service) GetResultsByQuizID(ctx context.Context, quizID int) ([]models.UsersResult, error) {
	ctx, span := s.tracer.Start(ctx, "resultService.GetResultsByQuizID")
	defer span.End()

	return s.repo.GetByQuizID(ctx, quizID)
}

func (s *Service) SubmitResult(ctx context.Context, userID, quizID int, input domain.SubmitResult) (models.UsersResult, error) {
	ctx, span := s.tracer.Start(ctx, "resultService.SubmitResult")
	defer span.End()

	if _, err := s.quizRepo.GetByID(ctx, quizID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UsersResult{}, err
	}

	return s.repo.SubmitResult(ctx, userID, input.AttemptID)
}
