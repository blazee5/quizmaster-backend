package repository

import (
	"context"
	"encoding/json"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/blazee5/testhub-backend/internal/quiz"
	"github.com/redis/go-redis/v9"
	"time"
)

type quizRedisRepo struct {
	redisClient *redis.Client
}

func NewAuthRedisRepo(redisClient *redis.Client) quiz.RedisRepository {
	return &quizRedisRepo{redisClient: redisClient}
}

func (repo *quizRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Quiz, error) {
	quizBytes, err := repo.redisClient.Get(ctx, "quiz:"+key).Bytes()

	if err != nil {
		return nil, err
	}

	var quiz *models.Quiz

	if err = json.Unmarshal(quizBytes, quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (repo *quizRedisRepo) SetQuizCtx(ctx context.Context, key string, seconds int, quiz *models.Quiz) error {
	quizBytes, err := json.Marshal(quiz)

	if err != nil {
		return err
	}

	if err := repo.redisClient.Set(ctx, "quiz:"+key, quizBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *quizRedisRepo) DeleteQuizCtx(ctx context.Context, key string) error {
	if err := repo.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
