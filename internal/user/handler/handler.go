package handler

import (
	"database/sql"
	"errors"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/user"
	"github.com/blazee5/testhub-backend/lib/http_utils"
	"github.com/blazee5/testhub-backend/lib/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	log     *zap.SugaredLogger
	service user.Service
}

func NewHandler(log *zap.SugaredLogger, service user.Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) Get(c echo.Context) error {
	userId := c.Get("userId").(int)

	user, err := h.service.GetById(c.Request().Context(), userId)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, domain.User{
		Id:     user.Id,
		Fio:    user.Fio,
		Email:  user.Email,
		Avatar: user.Avatar,
	})
}

func (h *Handler) GetQuizzes(c echo.Context) error {
	userId := c.Get("userId").(int)

	quizzes, err := h.service.GetQuizzes(c.Request().Context(), userId)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quizzes not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user quizzes: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) GetResults(c echo.Context) error {
	userId := c.Get("userId").(int)

	quizzes, err := h.service.GetResults(c.Request().Context(), userId)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "quizzes not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user results: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) Update(c echo.Context) error {
	var input domain.UpdateUser

	userId := c.Get("userId").(int)

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

	err := h.service.Update(c.Request().Context(), userId, input)

	if err != nil {
		h.log.Infof("error while update user: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "success")
}

func (h *Handler) Delete(c echo.Context) error {
	userId := c.Get("userId").(int)

	err := h.service.Delete(c.Request().Context(), userId)

	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user not found",
		})
	}

	if err != nil {
		h.log.Infof("error while get user results: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "OK")
}

func (h *Handler) UploadAvatar(c echo.Context) error {
	userId := c.Get("userId").(int)

	file, err := c.FormFile("file")

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "file is required",
		})
	}

	if err := http_utils.UploadFile(file, "./public/"+file.Filename); err != nil {
		return c.JSON(http.StatusInternalServerError, "error while upload avatar")
	}

	err = h.service.ChangeAvatar(c.Request().Context(), userId, file.Filename)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.String(http.StatusOK, "success")
}
