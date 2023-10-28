package main

import (
	"github.com/blazee5/testhub-backend/internal/postgres"
	"github.com/blazee5/testhub-backend/internal/routes"
	"github.com/blazee5/testhub-backend/lib/logger"
	libValidator "github.com/blazee5/testhub-backend/lib/validator"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log := logger.NewLogger()
	db := postgres.New()

	e := echo.New()
	e.Use(middleware.Recover())

	e.Validator = libValidator.NewValidator(validator.New())

	server := routes.NewServer(log, db)
	server.InitRoutes(e)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
