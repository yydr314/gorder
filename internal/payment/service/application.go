package service

import (
	"context"
	grpcClient "github.com/lingjun0314/goder/common/client"
	"github.com/lingjun0314/goder/common/metrics"
	"github.com/lingjun0314/goder/payment/adapters"
	"github.com/lingjun0314/goder/payment/app"
	"github.com/lingjun0314/goder/payment/app/command"
	"github.com/lingjun0314/goder/payment/domain"
	"github.com/lingjun0314/goder/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewApplication 膠水層，把所有要用的邏輯都進行依賴注入
func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	//memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}

}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
