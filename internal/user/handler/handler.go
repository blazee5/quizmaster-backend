package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	userService "github.com/blazee5/quizmaster-backend/internal/user"
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
	service userService.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service userService.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

func (h *Handler) GetMe(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "user.GetMe")
	defer span.End()

	userID := c.Get("userID").(int)

	user, err := h.service.GetByID(ctx, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetByID(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "user.GetByID")
	defer span.End()

	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	user, err := h.service.GetByID(ctx, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user by id: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Update(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "user.Update")
	defer span.End()

	var input domain.UpdateUser

	userID := c.Get("userID").(int)

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "bad request",
		})
	}

	if err := c.Validate(&input); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": response.ValidationError(validateErr),
		})
	}

	err := h.service.Update(ctx, userID, input)

	if err != nil {
		h.log.Infof("error while update user: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "success")
}

func (h *Handler) Delete(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "user.Delete")
	defer span.End()

	userID := c.Get("userID").(int)

	err := h.service.Delete(ctx, userID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user results: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) UploadAvatar(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "user.UploadAvatar")
	defer span.End()

	userID := c.Get("userID").(int)

	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "file is required",
		})
	}

	err = h.service.ChangeAvatar(ctx, userID, file)

	if errors.Is(err, http_errors.ErrInvalidImage) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid image",
		})
	}

	if err != nil {
		h.log.Infof("error while user upload avatar: %v", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "success")
}
