package handler

import (
	authRepo "github.com/blazee5/quizmaster-backend/internal/auth/repository"
	authService "github.com/blazee5/quizmaster-backend/internal/auth/service"
	"github.com/blazee5/quizmaster-backend/internal/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitAuthRoutes(authGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, tracer trace.Tracer) {
	repos := authRepo.NewRepository(db, tracer)
	services := authService.NewService(log, repos, tracer)
	handlers := NewHandler(log, services, tracer)

	authGroup.POST("/signup", handlers.SignUp)
	authGroup.POST("/signin", handlers.SignIn)
	authGroup.POST("/signout", handlers.SignOut, middleware.AuthMiddleware)
}
