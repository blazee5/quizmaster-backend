package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	quizService "github.com/blazee5/quizmaster-backend/internal/quiz"
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
	service quizService.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service quizService.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

// @Summary Get all quizzes
// @Tags quiz
// @Description Get all quizzes
// @ID get-all-quizzes
// @Accept json
// @Produce json
// @Param title query string false "title"
// @Param sortBy query string false "sortBy"
// @Param sortDir query string false "sortDir"
// @Param size query int false "size"
// @Param page query int false "page"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /quiz [get]
func (h *Handler) GetAllQuizzes(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.GetAllQuizzes")
	defer span.End()

	title := c.QueryParam("title")
	sortBy := c.QueryParam("sortBy")
	sortDir := c.QueryParam("sortDir")

	page, err := strconv.Atoi(c.QueryParam("page"))

	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))

	if err != nil || size < 1 {
		size = 10
	}

	quizzes, err := h.service.GetAll(ctx, title, sortBy, sortDir, page, size)

	if err != nil {
		h.log.Infof("error while get all quizzes: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quizzes)
}

// @Summary Get quiz
// @Tags quiz
// @Description Get quiz by id
// @ID get-quiz
// @Accept json
// @Produce json
// @Param id path int true "quizID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /quiz/{id} [get]
func (h *Handler) GetQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.GetQuiz")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	quiz, err := h.service.GetByID(ctx, id)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get quiz: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quiz)
}

// @Summary Create quiz
// @Tags quiz
// @Description Create quiz
// @ID create-quiz
// @Accept json
// @Produce json
// @Authorization BearerAuth "Authorization"
// @Param quiz body domain.Quiz true "Quiz"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /quiz [post]
func (h *Handler) CreateQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.CreateQuiz")
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

	id, err := h.service.Create(ctx, userID, input)

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

// @Summary Update quiz
// @Tags quiz
// @Description Update quiz
// @ID update-quiz
// @Accept json
// @Produce json
// @Authorization BearerAuth "Authorization"
// @Param quiz body domain.Quiz true "Quiz"
// @Param id path int true "quizID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /quiz/{id} [put]
func (h *Handler) UpdateQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.UpdateQuiz")
	defer span.End()

	var input domain.Quiz

	id, err := strconv.Atoi(c.Param("id"))
	userID := c.Get("userID").(int)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
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

	err = h.service.Update(ctx, userID, id, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if errors.Is(err, http_errors.ErrPermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while update quiz: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

// @Summary Delete quiz
// @Tags quiz
// @Description Delete quiz
// @ID delete-quiz
// @Accept json
// @Produce json
// @Authorization BearerAuth "Authorization"
// @Param id path int true "quizID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /quiz/{id} [delete]
func (h *Handler) DeleteQuiz(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.DeleteQuiz")
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))
	userID := c.Get("userID").(int)

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	err = h.service.Delete(ctx, userID, id)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if errors.Is(err, http_errors.ErrPermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while delete quiz: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "success")
}

// @Summary Upload image
// @Tags quiz
// @Description Upload image
// @ID upload-image
// @Accept json
// @Produce json
// @Authorization BearerAuth "Authorization"
// @Param id path int true "quizID"
// @Param image formData file true "image"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /quiz/{id}/image [post]
func (h *Handler) UploadImage(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.UploadImage")
	defer span.End()

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	file, err := c.FormFile("image")

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "image is required",
		})
	}

	err = h.service.UploadImage(ctx, userID, quizID, file)

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

	if errors.Is(err, http_errors.ErrPermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while upload quiz image: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

// @Summary Delete image
// @Tags quiz
// @Description Delete image
// @ID delete-image
// @Accept json
// @Produce json
// @Authorization BearerAuth "Authorization"
// @Param id path int true "quizID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 403 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /quiz/{id}/image [delete]
func (h *Handler) DeleteImage(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.DeleteImage")
	defer span.End()

	userID := c.Get("userID").(int)
	quizID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid quiz id",
		})
	}

	err = h.service.DeleteImage(ctx, userID, quizID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quiz not found",
		})
	}

	if errors.Is(err, http_errors.ErrPermissionDenied) {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "permission denied",
		})
	}

	if err != nil {
		h.log.Infof("error while delete quiz image: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
