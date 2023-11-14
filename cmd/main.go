package main

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/routes"
	"github.com/blazee5/quizmaster-backend/lib/db/postgres"
	"github.com/blazee5/quizmaster-backend/lib/db/redis"
	"github.com/blazee5/quizmaster-backend/lib/logger"
	libValidator "github.com/blazee5/quizmaster-backend/lib/validator"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	log := logger.NewLogger()
	db := postgres.New()
	rdb := redis.NewRedisClient()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Validator = libValidator.NewValidator(validator.New())
	server := routes.NewServer(log, db, rdb)
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
