package broker

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"time"
)

const (
	DLX                = "dlx"
	DLQ                = "dlq"
	amqpRetryHeaderKey = "x-amqp-retry"
)

var (
	maxRetryCount = viper.GetInt64("rabbitmq.max-retry")
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
	if err = createDLX(ch); err != nil {
		logrus.Fatal(err)
	}
	return ch, conn.Close
}

func createDLX(ch *amqp091.Channel) error {

	q, err := ch.QueueDeclare("share_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q.Name, "", DLX, false, nil)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)
	return err
}

func HandleRetry(ctx context.Context, ch *amqp091.Channel, d *amqp091.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp091.Table{}
	}
	retryCount, ok := d.Headers[amqpRetryHeaderKey].(int64)
	if !ok {
		retryCount = 0
	}
	retryCount++
	d.Headers[amqpRetryHeaderKey] = retryCount
	if retryCount >= maxRetryCount {
		logrus.Infof("moving message %s to  dlq", d.MessageId)
		return ch.PublishWithContext(ctx, "", DLQ, false, false, amqp091.Publishing{
			Headers:      d.Headers,
			ContentType:  "application/json",
			Body:         d.Body,
			DeliveryMode: amqp091.Persistent,
		})
	}

	logrus.Infof("retring message %s count=%d", d.MessageId, retryCount)
	time.Sleep(time.Second * time.Duration(retryCount))
	return ch.PublishWithContext(ctx, d.Exchange, d.RoutingKey, false, false, amqp091.Publishing{
		Headers:      d.Headers,
		ContentType:  "application/json",
		Body:         d.Body,
		DeliveryMode: amqp091.Persistent,
	})
}

type rabbitMQHeaderCarrier map[string]interface{}

// Get 實現 propagation 的 interface
func (r rabbitMQHeaderCarrier) Get(key string) string {
	value, ok := r[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (r rabbitMQHeaderCarrier) Set(key string, value string) {
	r[key] = value
}

func (r rabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, len(r))
	i := 0
	for key := range r {
		keys[i] = key
		i++
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(rabbitMQHeaderCarrier)
	// 使用 inject 方法將 context 內容放入 carrier 中
	otel.GetTextMapPropagator().Inject(ctx, &carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, rabbitMQHeaderCarrier(headers)) // 這裡的意思是把 headers 轉換為類型 RabbitMQHeaderCarrier
}
