package handler

import (
	"database/sql"
	"errors"
	adminauth "github.com/blazee5/quizmaster-backend/internal/admin/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service adminauth.Service
}

func NewHandler(log *zap.SugaredLogger, service adminauth.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) SignInAdmin(c echo.Context) error {
	var input domain.SignInRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	if err := c.Validate(&input); err != nil {
		validateErr := err.(validator.ValidationErrors)

		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": response.ValidationError(validateErr),
		})
	}

	token, err := h.service.GenerateToken(c.Request().Context(), input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "invalid credentials",
		})
	}

	if err != nil {
		h.log.Infof("error while sign in admin: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	cookie := auth.GenerateNewTokenCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
