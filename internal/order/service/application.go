package service

import (
	"context"
	grpcClient "github.com/lingjun0314/goder/common/client"
	"github.com/lingjun0314/goder/order/adapters/grpc"

	"github.com/lingjun0314/goder/common/metrics"
	"github.com/lingjun0314/goder/order/adapters"
	"github.com/lingjun0314/goder/order/app"
	"github.com/lingjun0314/goder/order/app/command"
	"github.com/lingjun0314/goder/order/app/query"
	"github.com/sirupsen/logrus"
)

// NewApplication 膠水層，把所有要用的邏輯都進行依賴注入
func NewApplication(ctx context.Context) (*app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	stockGRPC := grpc.NewStockGRPC(stockClient)
	if err != nil {
		panic(err)
	}
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}

}

func newApplication(_ context.Context, stockGRPC query.StockService) *app.Application {
	//	指定要用的依賴
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return &app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			//	在這裡依賴注入
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
