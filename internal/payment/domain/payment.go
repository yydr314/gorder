package domain

import (
	"context"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}
