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

func (q *quizRedisRepo) GetByIdCtx(ctx context.Context, key string) (*models.Quiz, error) {
	quizBytes, err := q.redisClient.Get(ctx, "quiz:"+key).Bytes()

	if err != nil {
		return nil, err
	}

	quiz := &models.Quiz{}

	if err = json.Unmarshal(quizBytes, quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (q *quizRedisRepo) SetQuizCtx(ctx context.Context, key string, seconds int, quiz *models.Quiz) error {
	quizBytes, err := json.Marshal(quiz)
	if err != nil {
		return err
	}

	if err := q.redisClient.Set(ctx, "quiz:"+key, quizBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (q *quizRedisRepo) DeleteQuizCtx(ctx context.Context, key string) error {
	if err := q.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}
