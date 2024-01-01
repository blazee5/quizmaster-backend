package routes

import (
	"context"
	adminAuthHandler "github.com/blazee5/quizmaster-backend/internal/admin/auth/handler"
	adminQuizHandler "github.com/blazee5/quizmaster-backend/internal/admin/quiz/handler"
	adminUserHandler "github.com/blazee5/quizmaster-backend/internal/admin/user/handler"
	answerHandler "github.com/blazee5/quizmaster-backend/internal/answer/handler"
	authHandler "github.com/blazee5/quizmaster-backend/internal/auth/handler"
	"github.com/blazee5/quizmaster-backend/internal/middleware"
	questionHandler "github.com/blazee5/quizmaster-backend/internal/question/handler"
	quizHandler "github.com/blazee5/quizmaster-backend/internal/quiz/handler"
	resultHandler "github.com/blazee5/quizmaster-backend/internal/result/handler"
	userHandler "github.com/blazee5/quizmaster-backend/internal/user/handler"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	socketio "github.com/vchitai/go-socket.io/v4"
)

func (s *Server) InitRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api")
	quizGroup := e.Group("/quiz")
	authGroup := e.Group("/auth")
	userGroup := apiGroup.Group("/user", middleware.AuthMiddleware)
	questionGroup := quizGroup.Group("/:id/questions", middleware.AuthMiddleware)
	answerGroup := questionGroup.Group("/:questionID/answers")
	adminGroup := e.Group("/admin")
	adminAuthGroup := adminGroup.Group("/auth")
	adminUsersGroup := adminGroup.Group("/users", middleware.AdminMiddleware)
	adminQuizzesGroup := adminGroup.Group("/quizzes", middleware.AdminMiddleware)

	authHandler.InitAuthRoutes(authGroup, s.log, s.db, s.rabbitConn, s.tracer)
	userHandler.InitUserRoutes(userGroup, s.log, s.db, s.rdb, s.awsClient, s.tracer)
	quizHandler.InitQuizRoutes(quizGroup, s.log, s.db, s.rdb, s.esClient, s.awsClient, s.tracer)
	resultHandler.InitResultRoutes(quizGroup, s.log, s.db, s.ws, s.tracer)
	questionHandler.InitQuestionRoutes(questionGroup, s.log, s.db, s.awsClient, s.tracer)
	answerHandler.InitAnswerRoutes(answerGroup, s.log, s.db, s.tracer)
	adminAuthHandler.InitAdminAuthRoutes(adminAuthGroup, s.log, s.db, s.tracer)
	adminUserHandler.InitAdminUserRoutes(adminUsersGroup, s.log, s.db, s.tracer)
	adminQuizHandler.InitAdminQuizRoutes(adminQuizzesGroup, s.log, s.db, s.tracer)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Any("/socket.io/", func(c echo.Context) error {
		s.ws.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	s.ws.OnConnect("/", func(s socketio.Conn, msg map[string]interface{}) error {
		s.SetContext(context.Background())

		return nil
	})
}
