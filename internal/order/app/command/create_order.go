package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lingjun0314/goder/common/broker"
	"github.com/lingjun0314/goder/common/decorator"
	"github.com/lingjun0314/goder/order/app/query"
	"github.com/lingjun0314/goder/order/convertor"
	domain "github.com/lingjun0314/goder/order/domain/order"
	"github.com/lingjun0314/goder/order/entity"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp091.Channel
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	channel *amqp091.Channel,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	if stockGRPC == nil {
		panic("nil stockGRPC")
	}
	if channel == nil {
		panic("nil channel")
	}
	return decorator.ApplyCommandDecorators(
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: stockGRPC,
			channel:   channel,
		},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// 聲明 channel
	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
	defer span.End()

	validateItems, err := c.validate(ctx, cmd.Items)

	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validateItems,
	})
	if err != nil {
		return nil, err
	}

	// 發送 MQ 消息
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	header := broker.InjectRabbitMQHeaders(ctx)
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Body:         marshalledOrder,
		Headers:      header,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

// 確認是否有庫存
func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWihhQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, err
	}
	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
}

// 如果請求的 items 有重複，那麼就要利用這個函式將數量合在一起
func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)

	//	綜合數量
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}

	var res []*entity.ItemWithQuantity
	// 回傳結果
	for id, quantity := range merged {
		res = append(res, &entity.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return res
}
