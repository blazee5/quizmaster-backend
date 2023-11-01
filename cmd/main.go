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
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	log := logger.NewLogger()
	db := postgres.New()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.Validator = libValidator.NewValidator(validator.New())

	server := routes.NewServer(log, db)
	server.InitRoutes(e)

	log.Fatal(e.Start(os.Getenv("PORT")))
}
