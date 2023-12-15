package handler

import (
	adminuser "github.com/blazee5/quizmaster-backend/internal/admin/user"
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
	service adminuser.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service adminuser.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

func (h *Handler) CreateUser(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.user.CreateUser")
	defer span.End()

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

	id, err := h.service.CreateUser(ctx, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		h.log.Infof("error while admin create user: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": id,
	})
}

func (h *Handler) GetUsers(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.user.GetUsers")
	defer span.End()

	users, err := h.service.GetUsers(ctx)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		h.log.Infof("error while admin get users: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.user.UpdateUser")
	defer span.End()

	var input domain.User

	userID, err := strconv.Atoi(c.Param("userID"))

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

	err = h.service.UpdateUser(ctx, userID, input)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		h.log.Infof("error while admin update user: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) DeleteUser(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.user.DeleteUser")
	defer span.End()

	userID, err := strconv.Atoi(c.Param("userID"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid id",
		})
	}

	err = h.service.DeleteUser(ctx, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		h.log.Infof("error while admin delete user: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
