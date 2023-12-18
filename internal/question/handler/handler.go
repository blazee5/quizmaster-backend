package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/question"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
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
	service question.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service question.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

// @Summary Create question
// @Tags question
// @Description create question
// @ID create-question
// @Accept json
// @Produce json
// @Param user body domain.Question true "question"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/question [post]
func (h *Handler) CreateQuestion(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.CreateQuestion")
	defer span.End()

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	id, err := h.service.Create(ctx, userID, quizID)

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
		h.log.Infof("error while create quiz: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) GetQuizQuestions(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.GetQuizQuestions")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	questions, err := h.service.GetQuestionsByID(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get questions by quiz id: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, questions)
}

func (h *Handler) UpdateQuestion(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.UpdateQuestion")
	defer span.End()

	var input domain.Question

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	questionID, err := strconv.Atoi(c.Param("questionID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	input.QuizID = quizID

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

	err = h.service.Update(ctx, questionID, userID, quizID, input)

	if err != nil {
		h.log.Infof("error while update question: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) DeleteQuestion(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.DeleteQuestion")
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

	err = h.service.Delete(ctx, questionID, userID, quizID)

	if errors.Is(err, http_errors.PermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while delete question: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) UploadImage(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.UploadImage")
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

	file, err := c.FormFile("image")

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "file is required",
		})
	}

	err = h.service.UploadImage(ctx, questionID, userID, quizID, file)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if errors.Is(err, http_errors.ErrInvalidImage) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid image",
		})
	}

	if errors.Is(err, http_errors.PermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while upload question image: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) DeleteImage(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.DeleteImage")
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

	err = h.service.DeleteImage(ctx, questionID, userID, quizID)

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
		h.log.Infof("error while delete question image: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) ChangeOrder(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "question.ChangeOrder")
	defer span.End()

	var input domain.ChangeQuestionOrder

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
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

	err = h.service.ChangeOrder(ctx, userID, quizID, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "question not found",
		})
	}

	if err != nil {
		h.log.Infof("error while change question order_id: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
