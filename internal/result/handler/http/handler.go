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
	"github.com/zishang520/socket.io/v2/socket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service result.Service
	ws      *socket.Server
}

func NewHandler(log *zap.SugaredLogger, service result.Service, ws *socket.Server) *Handler {
	return &Handler{log: log, service: service, ws: ws}
}

func (h *Handler) NewResult(c echo.Context) error {
	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	id, err := h.service.NewResult(c.Request().Context(), userID, quizID)

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

	err = h.service.SaveUserAnswer(c.Request().Context(), userID, quizID, input)

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

	result, err := h.service.SubmitResult(c.Request().Context(), userID, quizID, input)

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
	id := strconv.Itoa(quizID)

	results, err := h.service.GetResultsByQuizID(context.Background(), quizID)

	if errors.Is(err, sql.ErrNoRows) {
		return "quiz not found"
	}

	if err != nil {
		h.log.Infof("error while get quiz results: %s", err)
		return "server error"
	}

	h.ws.To(socket.Room("quiz:"+id)).Emit("message", results)

	return nil
}
