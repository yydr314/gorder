package integration

import (
	"context"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/price"
)

type StripeAPI struct {
	ApiKey string
}

func NewStripeAPI() *StripeAPI {
	return &StripeAPI{ApiKey: viper.GetString("stripe-key")}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	result, err := price.Get(pid, &stripe.PriceParams{})
	if err != nil {
		return "", err
	}

	return result.ID, err
}
