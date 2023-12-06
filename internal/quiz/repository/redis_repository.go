package repository

import (
	"context"
	"encoding/json"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type QuizRedisRepo struct {
	redisClient *redis.Client
	tracer      trace.Tracer
}

func NewQuizRedisRepo(redisClient *redis.Client, tracer trace.Tracer) *QuizRedisRepo {
	return &QuizRedisRepo{redisClient: redisClient, tracer: tracer}
}

func (repo *QuizRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "quizRedisRepo.GetByIDCtx")
	defer span.End()

	quizBytes, err := repo.redisClient.Get(ctx, "quiz:"+key).Bytes()

	if err != nil {
		return nil, err
	}

	var quiz *models.Quiz

	if err = json.Unmarshal(quizBytes, &quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (repo *QuizRedisRepo) SetQuizCtx(ctx context.Context, key string, seconds int, quiz *models.Quiz) error {
	ctx, span := repo.tracer.Start(ctx, "quizRedisRepo.SetQuizCtx")
	defer span.End()

	quizBytes, err := json.Marshal(quiz)

	if err != nil {
		return err
	}

	if err := repo.redisClient.Set(ctx, "quiz:"+key, quizBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *QuizRedisRepo) DeleteQuizCtx(ctx context.Context, key string) error {
	ctx, span := repo.tracer.Start(ctx, "quizRedisRepo.DeleteQuizCtx")
	defer span.End()

	if err := repo.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
