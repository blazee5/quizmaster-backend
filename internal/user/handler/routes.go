package handler

import (
	userRepo "github.com/blazee5/quizmaster-backend/internal/user/repository"
	userService "github.com/blazee5/quizmaster-backend/internal/user/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitUserRoutes(userGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, rdb *redis.Client, awsClient *minio.Client, tracer trace.Tracer) {
	repos := userRepo.NewRepository(db, tracer)
	redisRepos := userRepo.NewUserRedisRepo(rdb)
	awsRepos := userRepo.NewAWSRepository(awsClient)
	services := userService.NewService(log, repos, redisRepos, awsRepos, tracer)
	handlers := NewHandler(log, services, tracer)

	userGroup.GET("/me", handlers.GetMe)
	userGroup.GET("/:id", handlers.GetByID)
	userGroup.POST("/avatar", handlers.UploadAvatar)
	userGroup.PUT("", handlers.Update)
	userGroup.DELETE("", handlers.Delete)
}
