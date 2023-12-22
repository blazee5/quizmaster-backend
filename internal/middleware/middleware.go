package middleware

import (
	"github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Request().Cookie("token")

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "empty authorization cookie")
		}

		if token.Value == "" {
			return c.JSON(http.StatusUnauthorized, "empty authorization cookie")
		}

		userID, _, err := auth.ParseToken(token.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		c.Set("userID", userID)

		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := c.Request().Cookie("token")

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "empty authorization cookie")
		}

		if token.Value == "" {
			return c.JSON(http.StatusUnauthorized, "empty authorization cookie")
		}

		userID, roleID, err := auth.ParseToken(token.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		if roleID != 2 {
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": "forbidden",
			})
		}

		c.Set("userID", userID)

		return next(c)
	}
}
