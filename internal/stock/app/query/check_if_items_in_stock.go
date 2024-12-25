package query

import (
	"context"
	"github.com/lingjun0314/goder/common/decorator"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIdItemsInStockHandler struct {
	stockRepo stock.Repository
}

func NewCheckIfItemsInStockHandler(
	stockRepo stock.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators(
		checkIdItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

var stub = map[string]string{
	"1": "price_1QZurUEsREysJTgcLnc2lCox",
	"2": "price_4fewf1sd8fdse84sd",
}

func (c checkIdItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, i := range query.Items {
		priceID, ok := stub[i.ID]
		if !ok {
			priceID = stub["1"]
		}
		res = append(res, &orderpb.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}
