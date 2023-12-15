package handler

import (
	"database/sql"
	"errors"
	adminauth "github.com/blazee5/quizmaster-backend/internal/admin/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service adminauth.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service adminauth.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

func (h *Handler) SignInAdmin(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "admin.auth.SignInAdmin")
	defer span.End()

	var input domain.SignInRequest

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

	token, err := h.service.GenerateToken(ctx, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "invalid credentials",
		})
	}

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		h.log.Infof("error while sign in admin: %v", err)

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	cookie := auth.GenerateNewTokenCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
