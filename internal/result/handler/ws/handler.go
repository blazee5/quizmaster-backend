package ws

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/result"
	"github.com/labstack/gommon/log"
	"github.com/zishang520/socket.io/v2/socket"
	"go.uber.org/zap"
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

func (h *Handler) GetResults(msg ...any) {
	quizID := msg[0].(string)
	h.ws.SocketsJoin(socket.Room("quiz:" + quizID))

	id, err := strconv.Atoi(quizID)

	if err != nil {
		log.Info("invalid quizID")
		return
	}

	results, err := h.service.GetResultsByQuizID(context.Background(), id)

	if errors.Is(err, sql.ErrNoRows) {
		log.Info("quiz not found")
		return
	}

	if err != nil {
		h.log.Infof("error while get quiz results: %s", err)
		return
	}

	h.ws.Emit("message", results)
}
