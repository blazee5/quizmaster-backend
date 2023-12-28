package rabbitmq

import (
	"context"
	"github.com/blazee5/quizmaster-backend/lib/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
)

type Producer struct {
	log   *zap.SugaredLogger
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

func NewProducer(log *zap.SugaredLogger, conn *amqp.Connection) *Producer {
	return &Producer{log: log, conn: conn}
}

func (p *Producer) InitProducer() {
	ch, err := rabbitmq.NewChannelConn(p.conn)

	if err != nil {
		log.Fatalf("error while init producer: %v", err)
	}

	q, err := rabbitmq.NewQueueConn(ch)

	if err != nil {
		log.Fatalf("error while init producer: %v", err)
	}

	p.ch = ch
	p.queue = q
}

func (p *Producer) PublishMessage(ctx context.Context, msg []byte) error {
	err := p.ch.PublishWithContext(ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msg,
		})

	if err != nil {
		return err
	}

	return nil
}
