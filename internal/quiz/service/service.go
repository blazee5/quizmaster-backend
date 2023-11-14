package service

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/internal/user"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Service struct {
	log           *zap.SugaredLogger
	repo          quiz.Repository
	quizRedisRepo quiz.RedisRepository
	userRedisRepo user.RedisRepository
}

func NewService(log *zap.SugaredLogger, repo quiz.Repository, quizRedisRepo quiz.RedisRepository, userRedisRepo user.RedisRepository) *Service {
	return &Service{log: log, repo: repo, quizRedisRepo: quizRedisRepo, userRedisRepo: userRedisRepo}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Quiz, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetById(ctx context.Context, id int) (models.Quiz, error) {
	cachedQuiz, err := s.quizRedisRepo.GetByIdCtx(ctx, strconv.Itoa(id))

	if err != nil {
		s.log.Infof("cannot get quiz by id in redis: %v", err)
	}

	if cachedQuiz != nil {
		return *cachedQuiz, nil
	}

	quiz, err := s.repo.GetById(ctx, id)

	if err != nil {
		return models.Quiz{}, err
	}

	if err := s.quizRedisRepo.SetQuizCtx(ctx, strconv.Itoa(quiz.Id), 600, &quiz); err != nil {
		s.log.Infof("error while save quiz to cache: %v", err)
	}

	return quiz, nil
}

func (s *Service) Create(ctx context.Context, userId int, input domain.Quiz) (int, error) {
	id, err := s.repo.Create(ctx, userId, input)

	if err != nil {
		return 0, err
	}

	if err := s.userRedisRepo.DeleteUserCtx(ctx, strconv.Itoa(userId)); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) SaveResult(ctx context.Context, userId, quizId int, input domain.Result) (int, error) {
	totalQuestions, err := s.repo.GetQuestionsById(ctx, quizId, false)

	if err != nil {
		return 0, err
	}

	tx, err := s.repo.NewTx()

	if err != nil {
		return 0, err
	}

	score, err := s.SaveResultProcess(ctx, tx, userId, input)

	if err != nil {
		return 0, err
	}

	err = s.repo.SaveResult(ctx, userId, quizId, int(score), int(score/float64(len(totalQuestions))*100))

	if err != nil {
		return 0, err
	}

	if err := s.userRedisRepo.DeleteUserCtx(ctx, strconv.Itoa(userId)); err != nil {
		return 0, err
	}

	return int(score), nil
}

func (s *Service) Delete(ctx context.Context, userId, quizId int) error {
	quiz, err := s.repo.GetById(ctx, quizId)

	if err != nil {
		return err
	}

	if quiz.UserId != userId {
		return http_errors.PermissionDenied
	}

	return s.repo.Delete(ctx, quizId)
}

func (s *Service) SaveResultProcess(ctx context.Context, tx *sqlx.Tx, userId int, input domain.Result) (float64, error) {
	var score float64

	for questionId, answer := range input.Answers {
		var userCorrectAnswers int
		var totalUserAnswers int
		totalCorrectAnswers, err := s.repo.GetCorrectAnswers(ctx, questionId)

		if err != nil {
			return 0, err
		}

		questionType, err := s.repo.GetQuestionType(ctx, questionId)

		if err != nil {
			return 0, err
		}

		if questionType == "choice" {
			switch value := answer.(type) {
			case []interface{}:
				for _, answerId := range value {
					ans, err := s.repo.GetAnswerById(ctx, int(answerId.(float64)))

					if err != nil {
						return 0, err
					}

					err = s.repo.SaveUserAnswer(ctx, tx, userId, questionId, int(answerId.(float64)), "")

					if err != nil {
						return 0, err
					}

					if ans.IsCorrect {
						userCorrectAnswers++
					}

					totalUserAnswers++
				}
			}
		} else {
			switch value := answer.(type) {
			case string:
				answers, err := s.repo.GetAnswersById(ctx, questionId)

				if err != nil {
					return 0, err
				}

				err = s.repo.SaveUserAnswer(ctx, tx, userId, questionId, 0, value)

				if err != nil {
					return 0, err
				}

				for _, ans := range answers {
					if strings.ToLower(ans.Text) == value && ans.IsCorrect {
						userCorrectAnswers++
					}
				}

				totalUserAnswers++
			}
		}

		if userCorrectAnswers == totalCorrectAnswers && totalUserAnswers == userCorrectAnswers {
			score++
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return score, nil
}
