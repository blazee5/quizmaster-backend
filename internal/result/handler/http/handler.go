package http

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/result"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	socketio "github.com/vchitai/go-socket.io/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service result.Service
	ws      *socketio.Server
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service result.Service, ws *socketio.Server, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, ws: ws, tracer: tracer}
}

func (h *Handler) NewResult(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "result.NewResult")
	defer span.End()

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	id, err := h.service.NewResult(ctx, userID, quizID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while create result: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) SaveResult(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "result.SaveResult")
	defer span.End()

	var input domain.UserAnswer

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

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

	err = h.service.SaveUserAnswer(ctx, userID, quizID, input)

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
		h.log.Infof("error while save results: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) SubmitResult(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "result.SubmitResult")
	defer span.End()

	var input domain.SubmitResult

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

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

	result, err := h.service.SubmitResult(ctx, userID, quizID, input)

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
		h.log.Infof("error while submit result: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	h.UpdateResults(quizID)

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateResults(quizID int) interface{} {
	ctx, span := h.tracer.Start(context.Background(), "result.UpdateResults")
	defer span.End()

	id := strconv.Itoa(quizID)

	results, err := h.service.GetResultsByQuizID(ctx, quizID)

	if errors.Is(err, sql.ErrNoRows) {
		return "quiz not found"
	}

	if err != nil {
		h.log.Infof("error while get quiz results: %s", err)
		return "server error"
	}

	h.ws.BroadcastToRoom("/results", "quiz:"+id, "message", results)

	return nil
}
