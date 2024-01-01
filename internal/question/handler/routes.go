package handler

import (
	questionRepo "github.com/blazee5/quizmaster-backend/internal/question/repository"
	questionService "github.com/blazee5/quizmaster-backend/internal/question/service"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitQuestionRoutes(questionGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, awsClient *minio.Client, tracer trace.Tracer) {
	repos := questionRepo.NewRepository(db, tracer)
	awsRepos := questionRepo.NewAWSRepository(awsClient)
	quizRepos := quizRepo.NewRepository(db, tracer)
	services := questionService.NewService(log, repos, quizRepos, awsRepos, tracer)
	handlers := NewHandler(log, services, tracer)

	questionGroup.POST("", handlers.CreateQuestion)
	questionGroup.POST("/:questionID/image", handlers.UploadImage)
	questionGroup.GET("", handlers.GetQuizQuestions)
	questionGroup.GET("/author", handlers.GetQuestionsAuthor)
	questionGroup.PUT("/:questionID", handlers.UpdateQuestion)
	questionGroup.PUT("/order", handlers.ChangeOrder)
	questionGroup.DELETE("/:questionID", handlers.DeleteQuestion)
	questionGroup.DELETE("/:questionID/image", handlers.DeleteImage)
}
