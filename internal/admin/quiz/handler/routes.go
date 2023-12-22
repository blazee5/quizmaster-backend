package handler

import (
	adminQuizRepo "github.com/blazee5/quizmaster-backend/internal/admin/quiz/repository"
	adminQuizService "github.com/blazee5/quizmaster-backend/internal/admin/quiz/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitAdminQuizRoutes(adminQuizGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, tracer trace.Tracer) {
	repos := adminQuizRepo.NewRepository(db, tracer)
	services := adminQuizService.NewService(log, repos, tracer)
	handlers := NewHandler(log, services, tracer)

	adminQuizGroup.GET("", handlers.GetQuizzes)
	adminQuizGroup.POST("", handlers.CreateQuiz)
	adminQuizGroup.PUT("/:quizID", handlers.UpdateQuiz)
	adminQuizGroup.DELETE("/:quizID", handlers.DeleteQuiz)
}
