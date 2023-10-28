package handler

import (
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/user"
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

func (h *Handler) GetMe(c echo.Context) error {
	userId := c.Get("userId").(int)

	user, err := h.service.GetUserById(c.Request().Context(), userId)

	if err != nil {
		h.log.Infof("error while get user: %s", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "server error",
		})
	}

	return c.JSON(http.StatusOK, domain.User{
		Id:    user.Id,
		Fio:   user.Fio,
		Email: user.Email,
	})
}

func (h *Handler) UpdateMe(c echo.Context) error {
	return nil
}

func (h *Handler) DeleteMe(c echo.Context) error {
	return nil
}
