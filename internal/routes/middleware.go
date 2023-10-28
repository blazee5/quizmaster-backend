package routes

import (
	"github.com/blazee5/testhub-backend/lib/auth"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")

		if header == "" {
			return c.JSON(http.StatusUnauthorized, "empty authorization header")
		}

		headerParts := strings.Fields(header)
		if len(headerParts) != 2 {
			return c.JSON(http.StatusUnauthorized, "invalid authorization header")
		}

		userId, err := auth.ParseToken(headerParts[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set("userId", userId)

		return next(c)
	}
}
