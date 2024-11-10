package query

import (
	"context"
	"github.com/lingjun0314/goder/common/genproto/stockpb"

	"github.com/lingjun0314/goder/common/genproto/orderpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
