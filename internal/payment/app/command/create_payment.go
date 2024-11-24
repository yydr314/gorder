package command

import (
	"context"
	"github.com/lingjun0314/goder/common/decorator"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(processor domain.Processor, orderGRPC OrderService, logger *logrus.Entry, metricClient decorator.MetricsClient) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators(
		createPaymentHandler{processor: processor, orderGRPC: orderGRPC},
		logger,
		metricClient,
	)
}

func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	// 調用 stripe 服務獲取支付連結
	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for order: %s success, payment link: %s", cmd.Order.ID, link)
	// 建立新訂單
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}

	// 更新訂單狀態
	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	if err != nil {
		return "", err
	}

	return link, nil
}
