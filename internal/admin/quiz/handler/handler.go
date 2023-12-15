package handler

import (
	adminquiz "github.com/blazee5/quizmaster-backend/internal/admin/quiz"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service adminquiz.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service adminquiz.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

func (h *Handler) CreateQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.quiz.CreateQuiz")
	defer span.End()

	var input domain.Quiz
	userID := c.Get("userID").(int)

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

	id, err := h.service.CreateQuiz(ctx, userID, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		h.log.Infof("error while admin create quiz: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) GetQuizzes(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.quiz.GetQuizzes")
	defer span.End()

	quizzes, err := h.service.GetQuizzes(ctx)

	if err != nil {
		h.log.Infof("error while admin get quizzes: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) UpdateQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.quiz.UpdateQuiz")
	defer span.End()

	var input domain.Quiz

	quizID, err := strconv.Atoi(c.Param("quizID"))

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

	err = h.service.UpdateQuiz(ctx, quizID, input)

	if err != nil {
		h.log.Infof("error while admin update quiz")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) DeleteQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.quiz.DeleteQuiz")
	defer span.End()

	quizID, err := strconv.Atoi(c.Param("quizID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	err = h.service.DeleteQuiz(ctx, quizID)

	if err != nil {
		h.log.Infof("error while admin delete quiz: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
