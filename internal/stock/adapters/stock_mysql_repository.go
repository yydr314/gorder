package adapters

import (
	"context"
	"github.com/lingjun0314/goder/stock/entity"
	"github.com/lingjun0314/goder/stock/infrastructure/persistent"
)

type StockRepositoryMysql struct {
	db *persistent.MySQL
}

func NewStockRepositoryMysql(db *persistent.MySQL) *StockRepositoryMysql {
	return &StockRepositoryMysql{db: db}
}

func (s StockRepositoryMysql) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (s StockRepositoryMysql) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
	data, err := s.db.BatchGetStockByID(ctx, ids)
	if err != nil {
		return nil, err
	}

	var result []*entity.ItemWithQuantity
	for _, d := range data {
		result = append(result, &entity.ItemWithQuantity{
			ID:       d.ProductID,
			Quantity: d.Quantity,
		})
	}

	return result, nil
}
