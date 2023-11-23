package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service auth.Service
}

func NewHandler(log *zap.SugaredLogger, service auth.Service) *Handler {
	return &Handler{log: log, service: service}
}

// @Summary Sign up
// @Tags auth
// @Description Sign up
// @ID sign-up
// @Accept json
// @Produce json
// @Param user body domain.SignUpRequest true "user"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /auth/signup [post]
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

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "email already used",
			})
		}
	}

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

// @Summary Sign in
// @Tags auth
// @Description Sign in
// @ID sign-in
// @Accept json
// @Produce json
// @Param user body domain.SignInRequest true "user"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /auth/signin [post]
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

	cookie := authLib.GenerateNewTokenCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
