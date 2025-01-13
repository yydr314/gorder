package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lingjun0314/goder/common/broker"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"time"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, request *orderpb.Order) error
}

type Consumer struct {
	orderGRPC OrderService
}

func NewConsumer(orderGRPC OrderService) *Consumer {
	return &Consumer{orderGRPC: orderGRPC}
}

func (c *Consumer) Listen(ch *amqp091.Channel) {
	// 這裡的 exclusive 是 true 的原因是因為假設有多個 kitchen 服務，我們同一個 paid 消息只能被消費一次
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	if err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil); err != nil {
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
	var err error

	logrus.Infof("Kitchen  receive a message from %s, msg=%v", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))

	defer func() {
		span.End()

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
	if o.Status != "paid" {
		err = errors.New("order not paid, cannot cook")
		return
	}
	cook(o)

	span.AddEvent(fmt.Sprintf("order_cook: %v", o))
	o.Status = "ready"
	if err = c.orderGRPC.UpdateOrder(mqCtx, o); err != nil {
		if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
			logrus.Warnf("kitchen: error handling retry, err=%v", err)
		}
		return
	}

	span.AddEvent("kitchen.order.finished.updated")
	logrus.Info("consume success")
}

func cook(o *orderpb.Order) {
	logrus.Infof("cooking order: %s", o.ID)
	time.Sleep(5 * time.Second)
	logrus.Infof("order %s done!", o.ID)
}
