package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lingjun0314/goder/common/broker"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/payment/domain"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/webhook"
	"go.opentelemetry.io/otel"
	"io"
	"net/http"
)

type PaymentHandler struct {
	channel *amqp091.Channel
}

func NewPaymentHandler(ch *amqp091.Channel) *PaymentHandler {
	return &PaymentHandler{
		channel: ch,
	}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.HandleWebhook)
}

// 從 stripe doc 複製
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Infof("Error reading request body: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	// event 這裡要改，也是可以從 stripe doc 找到
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("endpoint-stripe-secret"))

	if err != nil {
		logrus.Infof("Error verifying webhook signature: %v\n", err)
		c.JSON(http.StatusBadRequest, err.Error()) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	// 關注完成支付
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			logrus.Infof("error unmarshal event.data.raw into session, err = %v", err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infof("payment for checkout session %v success!", session.ID)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				logrus.Infof("Error marshalling domain.order, err = %v", err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			tr := otel.Tracer("rabbitmq")
			mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderPaid))
			defer span.End()

			headers := broker.InjectRabbitMQHeaders(mqCtx)
			_ = h.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp091.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp091.Persistent,
				Body:         marshalledOrder,
				Headers:      headers,
			})
			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
		}
	}

	c.JSON(http.StatusOK, nil)
}
