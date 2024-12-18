package broker

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func Connect(user, password, host, port string) (*amqp091.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp091.Dial(address)
	if err != nil {
		logrus.Fatal(err)
	}
	// 開啟一個 channel
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return ch, conn.Close
}

type RabbitMQHeaderCarrier map[string]interface{}

// 實現 propagation 的 interface
func (r RabbitMQHeaderCarrier) Get(key string) string {
	value, ok := r[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (r RabbitMQHeaderCarrier) Set(key string, value string) {
	r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, len(r))
	i := 0
	for key := range r {
		keys[i] = key
		i++
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(RabbitMQHeaderCarrier)
	// 使用 inject 方法將 context 內容放入 carrier 中
	otel.GetTextMapPropagator().Inject(ctx, &carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers)) // 這裡的意思是把 headers 轉換為類型 RabbitMQHeaderCarrier
}
