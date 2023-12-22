package handler

import (
	answerRepo "github.com/blazee5/quizmaster-backend/internal/answer/repository"
	answerService "github.com/blazee5/quizmaster-backend/internal/answer/service"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitAnswerRoutes(answerGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, tracer trace.Tracer) {
	repos := answerRepo.NewRepository(db, tracer)
	quizRepos := quizRepo.NewRepository(db, tracer)
	services := answerService.NewService(log, repos, quizRepos, tracer)
	handlers := NewHandler(log, services, tracer)

	answerGroup.GET("", handlers.GetAnswers)
	answerGroup.POST("", handlers.CreateAnswer)
	answerGroup.PUT("/:answerID", handlers.UpdateAnswer)
	answerGroup.PUT("/order", handlers.ChangeOrder)
	answerGroup.DELETE("/:answerID", handlers.DeleteAnswer)
}
