package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/auth"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	authLib "github.com/blazee5/quizmaster-backend/lib/auth"
	"github.com/blazee5/quizmaster-backend/lib/http_errors"
	"github.com/blazee5/quizmaster-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service auth.Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service auth.Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

// @Summary Sign up
// @Tags auth
// @Description Sign up
// @ID sign-up
// @Accept json
// @Produce json
// @Param user body domain.SignUpRequest true "user"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /auth/signup [post]
func (h *Handler) SignUp(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "auth.SignUp")
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

	id, err := h.service.SignUp(ctx, input)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "email already used",
			})
		}
	}

	if err != nil {
		h.log.Infof("error while signup: %s", err)

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

// @Summary Sign in
// @Tags auth
// @Description Sign in
// @ID sign-in
// @Accept json
// @Produce json
// @Param user body domain.SignInRequest true "user"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /auth/signin [post]
func (h *Handler) SignIn(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "auth.SignIn")
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
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "invalid credentials",
		})
	}

	if err != nil {
		h.log.Infof("error while signin: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	cookie := authLib.GenerateNewTokenCookie(token)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

// @Summary Sign out
// @Tags auth
// @Description Sign out
// @ID sign-out
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /auth/signout [post]
func (h *Handler) SignOut(c echo.Context) error {
	_, span := h.tracer.Start(c.Request().Context(), "auth.SignOut")
	defer span.End()

	cookie := authLib.DeleteTokenCookie()
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

// @Summary Send Email Code
// @Tags auth
// @Description Send Code for reset email
// @ID send-code
// @Accept json
// @Produce json
// @Param code body domain.VerificationCode true "verification code"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /auth/send-email-code [post]
func (h *Handler) SendEmailCode(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "auth.SendEmailCode")
	defer span.End()

	var input domain.VerificationCode

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

	err := h.service.SendEmailCode(ctx, userID, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while send code on email: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

// @Summary Send Password Code
// @Tags auth
// @Description Send Code for reset password
// @ID send-code
// @Accept json
// @Produce json
// @Param code body domain.VerificationCode true "verification code"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /auth/send-password-code [post]
func (h *Handler) SendPasswordCode(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "auth.SendPasswordCode")
	defer span.End()

	var input domain.VerificationCode

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

	err := h.service.SendPasswordCode(ctx, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while send code on email: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

// @Summary Reset Email
// @Tags auth
// @Description Reset email
// @ID reset-email
// @Accept json
// @Produce json
// @Param email body domain.ResetEmailRequest true "code request"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /auth/reset-email [put]
func (h *Handler) ResetEmail(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "auth.ResetEmail")
	defer span.End()

	var input domain.ResetEmailRequest

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

	err := h.service.ResetEmail(ctx, userID, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "code not found",
		})
	}

	if errors.Is(err, http_errors.ErrCodeExpired) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "code is expired",
		})
	}

	if err != nil {
		h.log.Infof("error while reset email: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

// @Summary Reset Password
// @Tags auth
// @Description Reset password
// @ID reset-password
// @Accept json
// @Produce json
// @Param password body domain.ResetPasswordRequest true "new password and code"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /auth/reset-password [put]
func (h *Handler) ResetPassword(c echo.Context) error {
	ctx, span := h.tracer.Start(c.Request().Context(), "auth.ResetPassword")
	defer span.End()

	var input domain.ResetPasswordRequest

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

	err := h.service.ResetPassword(ctx, input)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "code not found",
		})
	}

	if errors.Is(err, http_errors.ErrCodeExpired) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "code is expired",
		})
	}

	if err != nil {
		h.log.Infof("error while reset password: %s", err)

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}
