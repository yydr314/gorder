package command

import (
	"context"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
