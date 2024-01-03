package handler

import (
	authRepo "github.com/blazee5/quizmaster-backend/internal/auth/repository"
	authService "github.com/blazee5/quizmaster-backend/internal/auth/service"
	"github.com/blazee5/quizmaster-backend/internal/middleware"
	"github.com/blazee5/quizmaster-backend/internal/rabbitmq"
	userRepo "github.com/blazee5/quizmaster-backend/internal/user/repository"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitAuthRoutes(authGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, rabbitConn *amqp.Connection, tracer trace.Tracer) {
	repos := authRepo.NewRepository(db, tracer)
	userRepos := userRepo.NewRepository(db, tracer)
	producer := rabbitmq.NewProducer(log, rabbitConn)
	producer.InitProducer()
	services := authService.NewService(log, repos, userRepos, producer, tracer)
	handlers := NewHandler(log, services, tracer)

	authGroup.POST("/signup", handlers.SignUp)
	authGroup.POST("/signin", handlers.SignIn)
	authGroup.POST("/signout", handlers.SignOut, middleware.AuthMiddleware)
	authGroup.POST("/send-code", handlers.SendCode)
	authGroup.PUT("/reset-email", handlers.ResetEmail, middleware.AuthMiddleware)
	authGroup.PUT("/reset-password", handlers.ResetPassword, middleware.AuthMiddleware)
}
