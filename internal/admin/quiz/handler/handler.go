package handler

import (
	adminquiz "github.com/blazee5/quizmaster-backend/internal/admin/quiz"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service adminquiz.Service
}

func NewHandler(log *zap.SugaredLogger, service adminquiz.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) CreateQuiz(c echo.Context) error {
	var input domain.Quiz
	userId := c.Get("userId").(int)

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

	id, err := h.service.CreateQuiz(c.Request().Context(), userId, input)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) GetQuizzes(c echo.Context) error {
	quizzes, err := h.service.GetQuizzes(c.Request().Context())

	if err != nil {
		h.log.Infof("error while update quiz: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) UpdateQuiz(c echo.Context) error {
	var input domain.Quiz

	quizId, err := strconv.Atoi(c.Param("quizId"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

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

	err = h.service.UpdateQuiz(c.Request().Context(), quizId, input)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) DeleteQuiz(c echo.Context) error {
	quizId, err := strconv.Atoi(c.Param("quizId"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	err = h.service.DeleteQuiz(c.Request().Context(), quizId)

	if err != nil {
		h.log.Infof("error while delete quiz: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
