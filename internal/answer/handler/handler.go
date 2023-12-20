package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/answer"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service answer.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service answer.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

func (h *Handler) CreateAnswer(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "answer.CreateAnswer")
	defer span.End()

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid question id",
		})
	}

	id, err := h.service.Create(ctx, userID, quizID, questionID)

	if errors.Is(err, http_errors.PermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while create quiz: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) GetAnswers(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "answer.GetAnswers")
	defer span.End()

	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid question id",
		})
	}

	answers, err := h.service.GetByQuestionID(ctx, quizID, questionID)

	if err != nil {
		h.log.Infof("error while get answers by question id: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, answers)
}

func (h *Handler) GetAnswersForUser(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "answer.GetAnswers")
	defer span.End()

	userID := c.Get("userID").(int)

	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid question id",
		})
	}

	answers, err := h.service.GetByQuestionIDForUser(ctx, quizID, questionID, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if errors.Is(err, http_errors.PermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while get answers by question id: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, answers)
}

func (h *Handler) UpdateAnswer(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "answer.UpdateAnswer")
	defer span.End()

	var input domain.Answer

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	answerID, err := strconv.Atoi(c.Param("answerID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid answer id",
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

	err = h.service.Update(ctx, answerID, userID, quizID, input)

	if err != nil {
		h.log.Infof("error while update answer: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) DeleteAnswer(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "answer.DeleteAnswer")
	defer span.End()

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	answerID, err := strconv.Atoi(c.Param("answerID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid answer id",
		})
	}

	err = h.service.Delete(ctx, answerID, userID, quizID)

	if errors.Is(err, http_errors.PermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while delete answer: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) ChangeOrder(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "answer.ChangeOrder")
	defer span.End()

	var input domain.ChangeAnswerOrder

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid question id",
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

	err = h.service.ChangeOrder(ctx, userID, quizID, questionID, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "question not found",
		})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
