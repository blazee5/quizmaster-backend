package main

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/routes"
	"github.com/blazee5/quizmaster-backend/lib/db/postgres"
	"github.com/blazee5/quizmaster-backend/lib/db/redis"
	"github.com/blazee5/quizmaster-backend/lib/elastic"
	"github.com/blazee5/quizmaster-backend/lib/logger"
	libValidator "github.com/blazee5/quizmaster-backend/lib/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/blazee5/quizmaster-backend/docs"
)

// @title QuizMaster Backend API
// @version 1.0
// @description This is a QuizMaster backend docs.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /

// @securitydefinitions.apikey ApiKeyAuth
// @in cookie
// @name token
// @tag.name auth
func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	panic("Error loading .env file")
	//}

	log := logger.NewLogger()
	db := postgres.New()
	rdb := redis.NewRedisClient()
	esclient := elastic.NewElasticSearchClient(log)

	e := echo.New()
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Validator = libValidator.NewValidator(validator.New())
	server := routes.NewServer(log, db, rdb, esclient)
	server.InitRoutes(e)

	log.Fatal(e.Start(os.Getenv("PORT")))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := e.Shutdown(context.Background()); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Infof("Error occured on db connection close: %v", err)
	}
}
