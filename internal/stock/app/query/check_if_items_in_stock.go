package query

import (
	"context"
	"github.com/lingjun0314/goder/common/decorator"
	"github.com/lingjun0314/goder/stock/domain/stock"
	"github.com/lingjun0314/goder/stock/entity"
	"github.com/lingjun0314/goder/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIdItemsInStockHandler struct {
	stockRepo stock.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(
	stockRepo stock.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	if stripeAPI == nil {
		panic("nil stripeAPI")
	}
	return decorator.ApplyQueryDecorators(
		checkIdItemsInStockHandler{
			stockRepo: stockRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricClient,
	)
}

func (c checkIdItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	if err := c.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}

	var res []*entity.Item
	for _, i := range query.Items {
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, i.ID)
		if err != nil {
			logrus.Warnf("GetPriceByProductID error, item ID=%s, err=%v", i.ID, err)
			continue
		}
		res = append(res, &entity.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}

func (c checkIdItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, i := range query {
		ids = append(ids, i.ID)
	}
	records, err := c.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}

	idQuantityMap := make(map[string]int32)

	for _, r := range records {
		idQuantityMap[r.ID] += r.Quantity
	}

	var failedOn []struct {
		ID   string
		Want int32
		Have int32
	}
	for _, q := range query {
		if idQuantityMap[q.ID] < q.Quantity {
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: q.ID, Want: q.Quantity, Have: idQuantityMap[q.ID]})
		}
	}

	if failedOn != nil {
		return stock.ExceedStockError{FailedOn: failedOn}
	}

	return nil
}
