package ws

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/result"
	socketio "github.com/vchitai/go-socket.io/v4"
	"go.uber.org/zap"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service result.Service
	ws      *socketio.Server
}

func NewHandler(log *zap.SugaredLogger, service result.Service, ws *socketio.Server) *Handler {
	return &Handler{log: log, service: service, ws: ws}
}

func (h *Handler) GetResults(conn socketio.Conn, quizID string) interface{} {
	conn.Join("quiz:" + quizID)

	id, err := strconv.Atoi(quizID)

	if err != nil {
		return "invalid quizID"
	}

	results, err := h.service.GetResultsByQuizID(context.Background(), id)

	if errors.Is(err, sql.ErrNoRows) {
		return "quiz not found"
	}

	if err != nil {
		h.log.Infof("error while get quiz results: %s", err)
		return "server error"
	}

	conn.Emit("message", results)
	return results
}
