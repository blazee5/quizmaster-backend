package rabbitmq

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/email"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"strconv"
)

type Consumer struct {
	log     *zap.SugaredLogger
	service email.Service
}

func NewConsumer(log *zap.SugaredLogger, service email.Service) *Consumer {
	return &Consumer{log: log, service: service}
}

func (c *Consumer) ConsumeQueue(ctx context.Context, ch *amqp.Channel) error {
	workers, err := strconv.Atoi(os.Getenv("WORKERS_COUNT"))

	if err != nil {
		c.log.Infof("invalid workers count in env: %v", err)

		return err
	}

	msgs, err := ch.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
		os.Getenv("RABBITMQ_CONSUMER"),
		false,
		false,
		false,
		false,
		nil,
	)

	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i <= workers; i++ {
		eg.Go(c.RunConsumer(ctx, msgs))
	}

	if err != nil {
		c.log.Infof("error in workers pool: %v", err)
	}

	return eg.Wait()
}

func (c *Consumer) RunConsumer(ctx context.Context, ch <-chan amqp.Delivery) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case msg, ok := <-ch:
				if !ok {
					c.log.Infof("channel is closed")
				}

				if err := c.service.SendEmail(string(msg.Body)); err != nil {
					c.log.Infof("error while send email: %v", err)
					continue
				}

				err := msg.Ack(false)

				if err != nil {
					c.log.Errorf("failed to acknowledge delivery: %v", err)
				}
			}
		}
	}
}
