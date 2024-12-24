package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lingjun0314/goder/common/broker"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/payment/app"
	"github.com/lingjun0314/goder/payment/app/command"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp091.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
	}

	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp091.Channel, msg amqp091.Delivery, q amqp091.Queue) {
	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	o := &orderpb.Order{}
	if err = json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmarshal msg to order, err=%v", err)
		return
	}
	if _, err = c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		logrus.Infof("failed to crate payment, err=%v", err)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("error handling retry, messageID=%s, err=%v", msg.MessageId, err)
		}
		return
	}

	span.AddEvent("payment.created")
	logrus.Info("consume success")
}
