package ws

import (
	"context"
	"database/sql"
	"errors"
	"github.com/blazee5/quizmaster-backend/internal/result"
	socketio "github.com/vchitai/go-socket.io/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"strconv"
)

type Handler struct {
	log     *zap.SugaredLogger
	service result.Service
	ws      *socketio.Server
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service result.Service, ws *socketio.Server, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, ws: ws, tracer: tracer}
}

func (h *Handler) GetResults(conn socketio.Conn, quizID string) interface{} {
	ctx, span := h.tracer.Start(context.Background(), "resultWs.GetResults")
	defer span.End()

	conn.Join("quiz:" + quizID)

	id, err := strconv.Atoi(quizID)

	if err != nil {
		return "invalid quizID"
	}

	results, err := h.service.GetResultsByQuizID(ctx, id)

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
