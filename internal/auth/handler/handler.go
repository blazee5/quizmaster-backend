package handler

import (
	"database/sql"
	"errors"
	auth "github.com/blazee5/testhub-backend/internal/auth/service"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service *auth.Service
}

func NewHandler(log *zap.SugaredLogger, service *auth.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) SignUp(c echo.Context) error {
	var input domain.SignUpRequest

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

	id, err := h.service.SignUp(c.Request().Context(), input)

	if err != nil {
		h.log.Infof("error while signup: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) SignIn(c echo.Context) error {
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
		h.log.Infof("error while signin: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
