package routes

import (
	authHandler "github.com/blazee5/testhub-backend/internal/auth/handler"
	authRepo "github.com/blazee5/testhub-backend/internal/auth/repository"
	authService "github.com/blazee5/testhub-backend/internal/auth/service"
	userHandler "github.com/blazee5/testhub-backend/internal/user/handler"
	userRepo "github.com/blazee5/testhub-backend/internal/user/repository"
	userService "github.com/blazee5/testhub-backend/internal/user/service"
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

	userRepos := userRepo.NewRepository(s.db)
	userServices := userService.NewService(userRepos)
	userHandlers := userHandler.NewHandler(s.log, userServices)

	user := e.Group("/user", AuthMiddleware)
	{
		user.GET("/me", userHandlers.GetMe)
		user.PUT("/", userHandlers.UpdateMe)
		user.DELETE("/", userHandlers.DeleteMe)
	}

	e.Static("/public", "public")
}
