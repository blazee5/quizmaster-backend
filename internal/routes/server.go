package routes

import (
	"context"
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
	resultHandler "github.com/blazee5/quizmaster-backend/internal/result/handler/http"
	"github.com/blazee5/quizmaster-backend/internal/result/handler/ws"
	resultRepo "github.com/blazee5/quizmaster-backend/internal/result/repository"
	resultService "github.com/blazee5/quizmaster-backend/internal/result/service"
	userHandler "github.com/blazee5/quizmaster-backend/internal/user/handler"
	userRepo "github.com/blazee5/quizmaster-backend/internal/user/repository"
	userService "github.com/blazee5/quizmaster-backend/internal/user/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
	socketio "github.com/vchitai/go-socket.io/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
	envProd        = "prod"
	certFile       = "/etc/letsencrypt/live/quizmaster.swedencentral.cloudapp.azure.com/fullchain.pem"
	keyFile        = "/etc/letsencrypt/live/quizmaster.swedencentral.cloudapp.azure.com/privkey.pem"
)

type Server struct {
	echo      *echo.Echo
	log       *zap.SugaredLogger
	db        *sqlx.DB
	rdb       *redis.Client
	esClient  *elasticsearch.Client
	ws        *socketio.Server
	tracer    trace.Tracer
	awsClient *minio.Client
}

func NewServer(echo *echo.Echo, log *zap.SugaredLogger, db *sqlx.DB, rdb *redis.Client, esClient *elasticsearch.Client, ws *socketio.Server, tracer trace.Tracer, awsClient *minio.Client) *Server {
	return &Server{echo: echo, log: log, db: db, rdb: rdb, esClient: esClient, ws: ws, tracer: tracer, awsClient: awsClient}
}

func (s *Server) Run() error {
	s.InitRoutes(s.echo)

	if os.Getenv("ENV") == envProd {
		go func() {
			s.log.Info("Server is listening on port 443")
			s.echo.Server.ReadTimeout = time.Second * 10
			s.echo.Server.WriteTimeout = time.Second * 10
			s.echo.Server.MaxHeaderBytes = maxHeaderBytes
			if err := s.echo.StartTLS(":443", certFile, keyFile); err != nil {
				s.log.Fatalf("Error starting TLS Server: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
		defer shutdown()

		s.log.Info("Server exiting...")
		return s.echo.Server.Shutdown(ctx)
	}

	go func() {
		if err := s.echo.Start(os.Getenv("PORT")); err != nil {
			s.log.Infof("Error Starting server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("Server exiting...")
	return s.echo.Server.Shutdown(ctx)
}

func (s *Server) InitRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	authRepos := authRepo.NewRepository(s.db, s.tracer)
	authServices := authService.NewService(s.log, authRepos, s.tracer)
	authHandlers := authHandler.NewHandler(s.log, authServices, s.tracer)

	auth := e.Group("/auth")
	{
		auth.POST("/signup", authHandlers.SignUp)
		auth.POST("/signin", authHandlers.SignIn)
	}

	api := e.Group("/api")
	{
		userRepos := userRepo.NewRepository(s.db, s.tracer)
		userRedisRepo := userRepo.NewUserRedisRepo(s.rdb)
		userAWSRepo := userRepo.NewAWSRepository(s.awsClient)
		userServices := userService.NewService(s.log, userRepos, userRedisRepo, userAWSRepo, s.tracer)
		userHandlers := userHandler.NewHandler(s.log, userServices, s.tracer)

		user := api.Group("/user", AuthMiddleware)
		{
			user.GET("/me", userHandlers.GetMe)
			user.GET("/:id", userHandlers.GetByID)
			user.POST("/avatar", userHandlers.UploadAvatar)
			user.PUT("", userHandlers.Update)
			user.DELETE("", userHandlers.Delete)
		}

		quizRepos := quizRepo.NewRepository(s.db, s.tracer)
		quizRedisRepo := quizRepo.NewQuizRedisRepo(s.rdb, s.tracer)
		quizElasticRepo := quizRepo.NewElasticRepository(s.esClient, s.tracer)
		quizAWSRepo := quizRepo.NewAWSRepository(s.awsClient)
		quizServices := quizService.NewService(s.log, quizRepos, quizRedisRepo, userRedisRepo, quizElasticRepo, quizAWSRepo, s.tracer)
		quizHandlers := quizHandler.NewHandler(s.log, quizServices, s.tracer)

		questionRepos := questionRepo.NewRepository(s.db, s.tracer)
		answerRepos := answerRepo.NewRepository(s.db, s.tracer)

		resultRepos := resultRepo.NewRepository(s.db, s.tracer)
		resultServices := resultService.NewService(s.log, resultRepos, quizRepos, questionRepos, answerRepos, s.tracer)
		resultHandlers := resultHandler.NewHandler(s.log, resultServices, s.ws, s.tracer)
		resultSocketHandlers := ws.NewHandler(s.log, resultServices, s.ws, s.tracer)

		s.ws.OnEvent("/results", "message", resultSocketHandlers.GetResults)

		quiz := e.Group("/quiz")
		{
			quiz.POST("", quizHandlers.CreateQuiz, AuthMiddleware)
			quiz.POST("/:id/image", quizHandlers.UploadImage, AuthMiddleware)
			quiz.GET("", quizHandlers.GetAllQuizzes)
			quiz.GET("/:id", quizHandlers.GetQuiz)
			quiz.PUT("/:id", quizHandlers.UpdateQuiz, AuthMiddleware)
			quiz.POST("/:id/start", resultHandlers.NewResult, AuthMiddleware)
			quiz.POST("/:id/save", resultHandlers.SaveResult, AuthMiddleware)
			quiz.POST("/:id/submit", resultHandlers.SubmitResult, AuthMiddleware)
			quiz.DELETE("/:id", quizHandlers.DeleteQuiz, AuthMiddleware)
			quiz.DELETE("/:id/image", quizHandlers.DeleteImage, AuthMiddleware)

			questionAWSRepos := questionRepo.NewAWSRepository(s.awsClient)
			questionServices := questionService.NewService(s.log, questionRepos, quizRepos, questionAWSRepos, s.tracer)
			questionHandlers := questionHandler.NewHandler(s.log, questionServices, s.tracer)

			question := quiz.Group("/:id/questions", AuthMiddleware)
			{
				question.POST("", questionHandlers.CreateQuestion)
				question.POST("/:questionID/image", questionHandlers.UploadImage)
				question.GET("", questionHandlers.GetQuizQuestions)
				question.PUT("/:questionID", questionHandlers.UpdateQuestion)
				question.PUT("/order", questionHandlers.ChangeOrder)
				question.DELETE("/:questionID", questionHandlers.DeleteQuestion)
				question.DELETE("/:questionID/image", questionHandlers.DeleteImage)

				answerServices := answerService.NewService(s.log, answerRepos, quizRepos, s.tracer)
				answerHandlers := answerHandler.NewHandler(s.log, answerServices, s.tracer)

				answer := question.Group("/:questionID/answers")
				{
					answer.GET("", answerHandlers.GetAnswers)
					answer.POST("", answerHandlers.CreateAnswer)
					answer.PUT("/:answerID", answerHandlers.UpdateAnswer)
					answer.PUT("/order", answerHandlers.ChangeOrder)
					answer.DELETE("/:answerID", answerHandlers.DeleteAnswer)
				}
			}
		}
	}

	admin := e.Group("/admin")
	{
		adminAuthRepos := adminAuthRepo.NewRepository(s.db, s.tracer)
		adminAuthServices := adminAuthService.NewService(s.log, adminAuthRepos, s.tracer)
		adminAuthHandlers := adminAuthHandler.NewHandler(s.log, adminAuthServices, s.tracer)

		auth := admin.Group("/auth")
		{
			auth.POST("/signin", adminAuthHandlers.SignInAdmin)
		}

		adminUserRepos := adminUserRepo.NewRepository(s.db, s.tracer)
		adminUserServices := adminUserService.NewService(s.log, adminUserRepos, s.tracer)
		adminUserHandlers := adminUserHandler.NewHandler(s.log, adminUserServices, s.tracer)

		users := admin.Group("/users", AdminMiddleware)
		{
			users.GET("", adminUserHandlers.GetUsers)
			users.POST("", adminUserHandlers.CreateUser)
			users.PUT("/:userID", adminUserHandlers.UpdateUser)
			users.DELETE("/:userID", adminUserHandlers.DeleteUser)
		}

		adminQuizRepos := adminQuizRepo.NewRepository(s.db, s.tracer)
		adminQuizServices := adminQuizService.NewService(s.log, adminQuizRepos, s.tracer)
		adminQuizHandlers := adminQuizHandler.NewHandler(s.log, adminQuizServices, s.tracer)

		quizzes := admin.Group("/quizzes", AdminMiddleware)
		{
			quizzes.GET("", adminQuizHandlers.GetQuizzes)
			quizzes.POST("", adminQuizHandlers.CreateQuiz)
			quizzes.PUT("/:quizID", adminQuizHandlers.UpdateQuiz)
			quizzes.DELETE("/:quizID", adminQuizHandlers.DeleteQuiz)
		}

		admin.GET("", adminUserHandlers.GetUsers)
	}

	e.Any("/socket.io/", func(c echo.Context) error {
		s.ws.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// websockets
	s.ws.OnConnect("/", func(s socketio.Conn, msg map[string]interface{}) error {
		s.SetContext(context.Background())

		return nil
	})
}
