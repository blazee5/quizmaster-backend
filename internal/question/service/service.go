package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	questionRepo "github.com/blazee5/quizmaster-backend/internal/question"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/files"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"mime/multipart"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     questionRepo.Repository
	quizRepo quizRepo.Repository
	awsRepo  questionRepo.AWSRepository
	tracer   trace.Tracer
}

func NewService(log *zap.SugaredLogger, repo questionRepo.Repository, quizRepo quizRepo.Repository, awsRepo questionRepo.AWSRepository, tracer trace.Tracer) *Service {
	return &Service{log: log, repo: repo, quizRepo: quizRepo, awsRepo: awsRepo, tracer: tracer}
}

func (s *Service) Create(ctx context.Context, userID, quizID int) (int, error) {
	ctx, span := s.tracer.Start(ctx, "questionService.Create")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	if quiz.UserID != userID {
		return 0, http_errors.ErrPermissionDenied
	}

	id, err := s.repo.CreateQuestion(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (s *Service) GetQuestionsByID(ctx context.Context, id int) ([]models.Question, error) {
	ctx, span := s.tracer.Start(ctx, "questionService.GetQuestionsByID")
	defer span.End()

	questions, err := s.repo.GetQuestionsByQuizID(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return questions, nil
}

func (s *Service) Test(ctx context.Context, id int) ([]models.QuestionWithAnswers, error) {
	ctx, span := s.tracer.Start(ctx, "questionService.Test")
	defer span.End()

	questions, err := s.repo.Test(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return questions, nil
}

func (s *Service) Update(ctx context.Context, id, userID, quizID int, input domain.Question) error {
	ctx, span := s.tracer.Start(ctx, "questionService.Update")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.Update(ctx, id, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id, userID, quizID int) error {
	ctx, span := s.tracer.Start(ctx, "questionService.Delete")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.Delete(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) UploadImage(ctx context.Context, id, userID, quizID int, fileHeader *multipart.FileHeader) error {
	ctx, span := s.tracer.Start(ctx, "questionService.UploadImage")
	defer span.End()

	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	question, err := s.repo.GetQuestionByID(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if quiz.UserID != userID || question.QuizID != quizID {
		return http_errors.ErrPermissionDenied
	}

	contentType, bytes, fileName, err := files.PrepareImage(fileHeader)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if question.Image != "" {
		err = s.awsRepo.DeleteFile(ctx, question.Image)

		if err != nil {
			return err
		}
	}

	err = s.awsRepo.SaveFile(ctx, fileName, contentType, bytes)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.UploadImage(ctx, id, fileName)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) DeleteImage(ctx context.Context, id, userID, quizID int) error {
	ctx, span := s.tracer.Start(ctx, "questionService.DeleteImage")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.DeleteImage(ctx, id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) ChangeOrder(ctx context.Context, userID, quizID int, input domain.QuestionOrder) error {
	ctx, span := s.tracer.Start(ctx, "questionService.ChangeOrder")
	defer span.End()

	err := s.checkPermissions(ctx, userID, quizID, input.QuestionID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	err = s.repo.ChangeOrder(ctx, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (s *Service) checkPermissions(ctx context.Context, userID, quizID, questionID int) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)

	if err != nil {
		return err
	}

	question, err := s.repo.GetQuestionByID(ctx, questionID)

	if err != nil {
		return err
	}

	if quiz.UserID != userID || question.QuizID != quizID {
		return http_errors.ErrPermissionDenied
	}

	return nil
}
