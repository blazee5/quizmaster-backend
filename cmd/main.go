package main

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/email/handler"
	"github.com/blazee5/quizmaster-backend/internal/routes"
	"github.com/blazee5/quizmaster-backend/lib/db/aws"
	"github.com/blazee5/quizmaster-backend/lib/db/postgres"
	"github.com/blazee5/quizmaster-backend/lib/db/redis"
	"github.com/blazee5/quizmaster-backend/lib/elastic"
	"github.com/blazee5/quizmaster-backend/lib/logger"
	"github.com/blazee5/quizmaster-backend/lib/rabbitmq"
	"github.com/blazee5/quizmaster-backend/lib/tracer"
	libValidator "github.com/blazee5/quizmaster-backend/lib/validator"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	socketio "github.com/vchitai/go-socket.io/v4"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/blazee5/quizmaster-backend/docs"
)

// @title QuizMaster Backend API
// @version 1.0
// @description This is a QuizMaster backend docs.
// @contact.name blaze
// @contact.url https://www.github.com/blazee5
// @host localhost:3000
// @BasePath /
// @securitydefinitions.apikey ApiKeyAuth
// @in cookie
// @name token
// @tag.name auth
func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	log := logger.NewLogger()
	db := postgres.New()
	rdb := redis.NewRedisClient()
	esClient := elastic.NewElasticSearchClient(log)
	ws := socketio.NewServer(nil)
	trace := tracer.InitTracer("Quizmaster")
	awsClient := aws.NewAWSClient()
	rabbitConn := rabbitmq.NewRabbitMQConn()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://quizer-opal.vercel.app", "https://quizmaster-admin.vercel.app"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	e.Validator = libValidator.NewValidator(validator.New())
	server := routes.NewServer(e, log, db, rdb, esClient, ws, trace, awsClient, rabbitConn)

	go func() {
		log.Fatal(server.Run())
	}()

	go func() {
		log.Fatal(ws.Serve())
	}()

	go func() {
		handler.InitEmailConsumer(context.Background(), log, rabbitConn)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := db.Close(); err != nil {
		log.Infof("Error occured on db connection close: %v", err)
	}

	if err := ws.Close(); err != nil {
		log.Infof("Error while socket server shutting down: %v", err)
	}
}
