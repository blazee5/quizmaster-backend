package handler

import (
	"github.com/blazee5/quizmaster-backend/internal/middleware"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz/repository"
	quizService "github.com/blazee5/quizmaster-backend/internal/quiz/service"
	"github.com/blazee5/quizmaster-backend/internal/user/repository"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitQuizRoutes(quizGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, rdb *redis.Client, esClient *elasticsearch.Client, awsClient *minio.Client, tracer trace.Tracer) {
	quizRepos := quizRepo.NewRepository(db, tracer)
	quizRedisRepos := quizRepo.NewQuizRedisRepo(rdb, tracer)
	quizElasticRepos := quizRepo.NewElasticRepository(esClient, tracer)
	quizAWSRepos := quizRepo.NewAWSRepository(awsClient)
	userRedisRepos := repository.NewUserRedisRepo(rdb)
	quizServices := quizService.NewService(log, quizRepos, quizRedisRepos, userRedisRepos, quizElasticRepos, quizAWSRepos, tracer)
	handlers := NewHandler(log, quizServices, tracer)

	quizGroup.POST("", handlers.CreateQuiz, middleware.AuthMiddleware)
	quizGroup.POST("/:id/image", handlers.UploadImage, middleware.AuthMiddleware)
	quizGroup.GET("", handlers.GetAllQuizzes)
	quizGroup.GET("/:id", handlers.GetQuiz)
	quizGroup.PUT("/:id", handlers.UpdateQuiz, middleware.AuthMiddleware)
	quizGroup.DELETE("/:id", handlers.DeleteQuiz, middleware.AuthMiddleware)
	quizGroup.DELETE("/:id/image", handlers.DeleteImage, middleware.AuthMiddleware)
}
