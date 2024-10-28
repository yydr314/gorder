package service

import (
	"context"

	"github.com/lingjun0314/goder/common/metrics"
	"github.com/lingjun0314/goder/order/adapters"
	"github.com/lingjun0314/goder/order/app"
	"github.com/lingjun0314/goder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) *app.Application {
	//	指定要用的依賴
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}

	return &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			//	在這裡依賴注入
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
