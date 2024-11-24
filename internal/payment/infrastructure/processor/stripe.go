package processor

import (
	"context"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	return &StripeProcessor{apiKey: apiKey}
}

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
