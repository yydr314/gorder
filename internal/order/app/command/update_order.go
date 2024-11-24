package command

import (
	"context"

	"github.com/lingjun0314/goder/common/decorator"
	domain "github.com/lingjun0314/goder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order    *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderResult struct {
	OrderID string
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	// stockGRPC
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) UpdateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyCommandDecorators(
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (u updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	if cmd.UpdateFn == nil {
		logrus.Warnf("updateOrderHandler got nil UpdateFn,order=%v", cmd.Order)
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	err := u.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
