package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/quiz"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/blazee5/quizmaster-backend/lib/http_utils"
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
	service quiz.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service quiz.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

func (h *Handler) GetAllQuizzes(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.GetAllQuizzes")
	defer span.End()

	quizzes, err := h.service.GetAll(ctx)

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

func (h *Handler) SearchByTitle(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "quiz.SearchByTitle")
	defer span.End()

	title := c.QueryParam("title")

	quizzes, err := h.service.Search(ctx, title)

	if err != nil {
		h.log.Infof("error while search quiz: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, quizzes)
}

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

	if errors.Is(err, http_errors.PermissionDenied) {
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

	if errors.Is(err, http_errors.PermissionDenied) {
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

	return c.String(http.StatusOK, "sucess")
}

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

	if err == nil {
		if err := http_utils.UploadFile(file, "public/"+file.Filename); err != nil {
			h.log.Infof("error while save question image: %v", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "server error",
			})
		}
	} else {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "image is required",
		})
	}

	err = h.service.UploadImage(ctx, userID, quizID, file.Filename)

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
		h.log.Infof("error while upload quiz image: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

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

	file, err := c.FormFile("image")

	if err == nil {
		if err := http_utils.UploadFile(file, "public/"+file.Filename); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "server error",
			})
		}
	}

	err = h.service.DeleteImage(ctx, userID, quizID)

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
		h.log.Infof("error while delete quiz image: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
