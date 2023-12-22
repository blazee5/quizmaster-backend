package handler

import (
	adminUserRepo "github.com/blazee5/quizmaster-backend/internal/admin/user/repository"
	adminUserService "github.com/blazee5/quizmaster-backend/internal/admin/user/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitAdminUserRoutes(adminUserGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, tracer trace.Tracer) {
	repos := adminUserRepo.NewRepository(db, tracer)
	services := adminUserService.NewService(log, repos, tracer)
	handlers := NewHandler(log, services, tracer)

	adminUserGroup.GET("", handlers.GetUsers)
	adminUserGroup.POST("", handlers.CreateUser)
	adminUserGroup.PUT("/:userID", handlers.UpdateUser)
	adminUserGroup.DELETE("/:userID", handlers.DeleteUser)
}
