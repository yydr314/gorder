package query

import (
	"context"

	"github.com/lingjun0314/goder/common/decorator"
	domain "github.com/lingjun0314/goder/order/domain/order"
	"github.com/sirupsen/logrus"
)

// 查詢條件的資料
type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

// 定義一個類型別名，用於暴露給外部使用
type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

// 查詢的具體實現結構體
type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

// Get customer order handler 的建構函式
func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyQueryDecorators(
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

// 查詢的具體實現
func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}
