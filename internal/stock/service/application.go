package service

import (
	"context"
	"github.com/lingjun0314/goder/common/metrics"
	"github.com/lingjun0314/goder/stock/adapters"
	"github.com/lingjun0314/goder/stock/app/query"
	"github.com/sirupsen/logrus"

	"github.com/lingjun0314/goder/stock/app"
)

func NewApplication(ctx context.Context) *app.Application {
	stockRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricClient),
		},
	}
}
