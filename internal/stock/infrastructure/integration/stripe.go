package integration

import (
	"context"
	_ "github.com/lingjun0314/goder/common/config"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/product"
)

type StripeAPI struct {
	ApiKey string
}

func NewStripeAPI() *StripeAPI {
	return &StripeAPI{ApiKey: viper.GetString("stripe-key")}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	stripe.Key = s.ApiKey
	params := &stripe.ProductParams{}
	result, err := product.Get(pid, params)
	if err != nil {
		return "", err
	}

	return result.DefaultPrice.ID, err
}
