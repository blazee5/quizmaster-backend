package routes

import (
	adminAuthHandler "github.com/blazee5/quizmaster-backend/internal/admin/auth/handler"
	adminAuthRepo "github.com/blazee5/quizmaster-backend/internal/admin/auth/repository"
	adminAuthService "github.com/blazee5/quizmaster-backend/internal/admin/auth/service"
	adminQuizHandler "github.com/blazee5/quizmaster-backend/internal/admin/quiz/handler"
	adminQuizRepo "github.com/blazee5/quizmaster-backend/internal/admin/quiz/repository"
	adminQuizService "github.com/blazee5/quizmaster-backend/internal/admin/quiz/service"
	adminUserHandler "github.com/blazee5/quizmaster-backend/internal/admin/user/handler"
	adminUserRepo "github.com/blazee5/quizmaster-backend/internal/admin/user/repository"
	adminUserService "github.com/blazee5/quizmaster-backend/internal/admin/user/service"
	answerHandler "github.com/blazee5/quizmaster-backend/internal/answer/handler"
	answerRepo "github.com/blazee5/quizmaster-backend/internal/answer/repository"
	answerService "github.com/blazee5/quizmaster-backend/internal/answer/service"
	authHandler "github.com/blazee5/quizmaster-backend/internal/auth/handler"
	authRepo "github.com/blazee5/quizmaster-backend/internal/auth/repository"
	authService "github.com/blazee5/quizmaster-backend/internal/auth/service"
	questionHandler "github.com/blazee5/quizmaster-backend/internal/question/handler"
	questionRepo "github.com/blazee5/quizmaster-backend/internal/question/repository"
	questionService "github.com/blazee5/quizmaster-backend/internal/question/service"
	quizHandler "github.com/blazee5/quizmaster-backend/internal/quiz/handler"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz/repository"
	quizService "github.com/blazee5/quizmaster-backend/internal/quiz/service"
	userHandler "github.com/blazee5/quizmaster-backend/internal/user/handler"
	userRepo "github.com/blazee5/quizmaster-backend/internal/user/repository"
	userService "github.com/blazee5/quizmaster-backend/internal/user/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Server struct {
	log      *zap.SugaredLogger
	db       *sqlx.DB
	rdb      *redis.Client
	esclient *elasticsearch.Client
}

func NewServer(log *zap.SugaredLogger, db *sqlx.DB, rdb *redis.Client, esclient *elasticsearch.Client) *Server {
	return &Server{log: log, db: db, rdb: rdb, esclient: esclient}
}

func (s *Server) InitRoutes(e *echo.Echo) {
	authRepos := authRepo.NewRepository(s.db)
	authServices := authService.NewService(s.log, authRepos)
	authHandlers := authHandler.NewHandler(s.log, authServices)

	auth := e.Group("/auth")
	{
		auth.POST("/signup", authHandlers.SignUp)
		auth.POST("/signin", authHandlers.SignIn)
	}

	api := e.Group("/api")
	{
		userRepos := userRepo.NewRepository(s.db)
		userRedisRepo := userRepo.NewUserRedisRepo(s.rdb)
		userServices := userService.NewService(s.log, userRepos, userRedisRepo)
		userHandlers := userHandler.NewHandler(s.log, userServices)

		user := api.Group("/user", AuthMiddleware)
		{
			user.GET("/me", userHandlers.Get)
			user.GET("/:id", userHandlers.GetByID)
			user.POST("/avatar", userHandlers.UploadAvatar)
			user.PUT("", userHandlers.Update)
			user.DELETE("", userHandlers.Delete)
		}

		quizRepos := quizRepo.NewRepository(s.db)
		quizRedisRepo := quizRepo.NewAuthRedisRepo(s.rdb)
		quizElasticRepo := quizRepo.NewElasticRepository(s.esclient)
		quizServices := quizService.NewService(s.log, quizRepos, quizRedisRepo, userRedisRepo, quizElasticRepo)
		quizHandlers := quizHandler.NewHandler(s.log, quizServices)

		quiz := e.Group("/quiz")
		{
			quiz.POST("", quizHandlers.CreateQuiz, AuthMiddleware)
			quiz.POST("/:id/image", quizHandlers.UploadImage, AuthMiddleware)
			quiz.GET("", quizHandlers.GetAllQuizzes)
			// quiz.GET("/search", quizHandlers.SearchByTitle)
			quiz.GET("/:id", quizHandlers.GetQuiz)
			quiz.PUT("/:id", quizHandlers.UpdateQuiz, AuthMiddleware)
			quiz.POST("/:id/save", quizHandlers.SaveResult, AuthMiddleware)
			quiz.DELETE("/:id", quizHandlers.DeleteQuiz, AuthMiddleware)
			quiz.DELETE("/:id/image", quizHandlers.DeleteImage, AuthMiddleware)

			questionRepos := questionRepo.NewRepository(s.db)
			questionServices := questionService.NewService(s.log, questionRepos, quizRepos)
			questionHandlers := questionHandler.NewHandler(s.log, questionServices)

			question := quiz.Group("/:id/questions", AuthMiddleware)
			{
				question.POST("", questionHandlers.CreateQuestion)
				question.POST("/:questionID/image", questionHandlers.UploadImage)
				question.GET("", questionHandlers.GetQuizQuestions)
				question.GET("/all", questionHandlers.GetAllQuizQuestions)
				question.PUT("/:questionID", questionHandlers.UpdateQuestion)
				question.DELETE("/:questionID", questionHandlers.DeleteQuestion)
				question.DELETE("/:questionID/image", questionHandlers.DeleteImage)

				answerRepos := answerRepo.NewRepository(s.db)
				answerServices := answerService.NewService(s.log, answerRepos, quizRepos)
				answerHandlers := answerHandler.NewHandler(s.log, answerServices)

				answer := question.Group("/:questionID/answers")
				{
					answer.POST("", answerHandlers.CreateAnswer)
					answer.PUT("/:answerID", answerHandlers.UpdateAnswer)
					answer.DELETE("/:answerID", answerHandlers.DeleteAnswer)
				}
			}
		}
	}

	admin := e.Group("/admin")
	{
		adminAuthRepos := adminAuthRepo.NewRepository(s.db)
		adminAuthServices := adminAuthService.NewService(s.log, adminAuthRepos)
		adminAuthHandlers := adminAuthHandler.NewHandler(s.log, adminAuthServices)

		auth := admin.Group("/auth")
		{
			auth.POST("/signin", adminAuthHandlers.SignInAdmin)
		}

		adminUserRepos := adminUserRepo.NewRepository(s.db)
		adminUserServices := adminUserService.NewService(s.log, adminUserRepos)
		adminUserHandlers := adminUserHandler.NewHandler(s.log, adminUserServices)

		users := admin.Group("/users", AdminMiddleware)
		{
			users.GET("", adminUserHandlers.GetUsers)
			users.POST("", adminUserHandlers.CreateUser)
			users.PUT("/:userID", adminUserHandlers.UpdateUser)
			users.DELETE("/:userID", adminUserHandlers.DeleteUser)
		}

		adminQuizRepos := adminQuizRepo.NewRepository(s.db)
		adminQuizServices := adminQuizService.NewService(s.log, adminQuizRepos)
		adminQuizHandlers := adminQuizHandler.NewHandler(s.log, adminQuizServices)

		quizzes := admin.Group("/quizzes", AdminMiddleware)
		{
			quizzes.GET("", adminQuizHandlers.GetQuizzes)
			quizzes.POST("", adminQuizHandlers.CreateQuiz)
			quizzes.PUT("/:quizID", adminQuizHandlers.UpdateQuiz)
			quizzes.DELETE("/:quizID", adminQuizHandlers.DeleteQuiz)
		}

		admin.GET("", adminUserHandlers.GetUsers)
	}

	e.Static("/public", "public")
}
