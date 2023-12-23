package handler

import (
	adminAuthRepo "github.com/blazee5/quizmaster-backend/internal/admin/auth/repository"
	adminAuthService "github.com/blazee5/quizmaster-backend/internal/admin/auth/service"
	"github.com/blazee5/quizmaster-backend/internal/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitAdminAuthRoutes(adminAuthGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, tracer trace.Tracer) {
	repos := adminAuthRepo.NewRepository(db, tracer)
	services := adminAuthService.NewService(log, repos, tracer)
	handlers := NewHandler(log, services, tracer)

	adminAuthGroup.POST("/signin", handlers.SignInAdmin)
	adminAuthGroup.POST("/signout", handlers.SignOutAdmin, middleware.AdminMiddleware)
}
