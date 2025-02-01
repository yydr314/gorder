package query

import (
	"context"
	"github.com/lingjun0314/goder/common/decorator"
	"github.com/lingjun0314/goder/stock/domain/stock"
	"github.com/lingjun0314/goder/stock/entity"
	"github.com/sirupsen/logrus"
)

// 查詢的結構資料
type GetItems struct {
	ItemIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]

type getItemsHandler struct {
	stockRepo stock.Repository
}

func NewGetItemsHandler(
	stockRepo stock.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetItemsHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators(
		getItemsHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*entity.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		return nil, err
	}
	return items, nil
}
