package routes

import (
	authHandler "github.com/blazee5/testhub-backend/internal/auth/handler"
	authRepo "github.com/blazee5/testhub-backend/internal/auth/repository"
	authService "github.com/blazee5/testhub-backend/internal/auth/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Server struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

func NewServer(log *zap.SugaredLogger, db *sqlx.DB) *Server {
	return &Server{log: log, db: db}
}

func (s *Server) InitRoutes(e *echo.Echo) {
	authRepos := authRepo.NewRepository(s.db)
	authServices := authService.NewService(authRepos)
	authHandlers := authHandler.NewHandler(s.log, authServices)

	auth := e.Group("/auth")
	{
		auth.POST("/signup", authHandlers.SignUp)
		auth.POST("/signin", authHandlers.SignIn)
	}
}
