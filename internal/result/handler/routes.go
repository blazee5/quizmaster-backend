package handler

import (
	answerRepo "github.com/blazee5/quizmaster-backend/internal/answer/repository"
	"github.com/blazee5/quizmaster-backend/internal/middleware"
	questionRepo "github.com/blazee5/quizmaster-backend/internal/question/repository"
	quizRepo "github.com/blazee5/quizmaster-backend/internal/quiz/repository"
	"github.com/blazee5/quizmaster-backend/internal/result/handler/http"
	wsHandler "github.com/blazee5/quizmaster-backend/internal/result/handler/ws"
	resultRepo "github.com/blazee5/quizmaster-backend/internal/result/repository"
	resultService "github.com/blazee5/quizmaster-backend/internal/result/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	socketio "github.com/vchitai/go-socket.io/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func InitResultRoutes(resultGroup *echo.Group, log *zap.SugaredLogger, db *sqlx.DB, ws *socketio.Server, tracer trace.Tracer) {
	repos := resultRepo.NewRepository(db, tracer)
	quizRepos := quizRepo.NewRepository(db, tracer)
	questionRepos := questionRepo.NewRepository(db, tracer)
	answerRepos := answerRepo.NewRepository(db, tracer)
	services := resultService.NewService(log, repos, quizRepos, questionRepos, answerRepos, tracer)
	handlers := http.NewHandler(log, services, ws, tracer)
	wsHandlers := wsHandler.NewHandler(log, services, ws, tracer)

	resultGroup.POST("/:id/start", handlers.NewResult, middleware.AuthMiddleware)
	resultGroup.POST("/:id/save", handlers.SaveResult, middleware.AuthMiddleware)
	resultGroup.POST("/:id/submit", handlers.SubmitResult, middleware.AuthMiddleware)

	ws.OnEvent("/results", "message", wsHandlers.GetResults)
}
