package handler

import (
	"context"
	rabbitmqHandler "github.com/blazee5/quizmaster-backend/internal/email/handler/rabbitmq"
	"github.com/blazee5/quizmaster-backend/internal/email/service"
	rabbitmq2 "github.com/blazee5/quizmaster-backend/lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func InitEmailConsumer(ctx context.Context, log *zap.SugaredLogger, rabbitConn *amqp.Connection) {
	emailServices := service.NewService(log)
	consumer := rabbitmqHandler.NewConsumer(log, emailServices)
	ch, err := rabbitmq2.NewChannelConn(rabbitConn)

	if err != nil {
		log.Fatalf("error while create channel in rabbitmq: %v", err)
	}

	err = consumer.ConsumeQueue(ctx, ch)

	if err != nil {
		log.Fatalf("error while consume queue: %v", err)
	}
}
