package integration

import "github.com/spf13/viper"

type StripeAPI struct {
	ApiKey string
}

func NewStripeAPI() *StripeAPI {
	return &StripeAPI{ApiKey: viper.GetString("stripe-key")}
}
