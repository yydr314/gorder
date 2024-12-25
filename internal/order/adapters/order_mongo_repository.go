package adapters

import (
	"context"
	domain "github.com/lingjun0314/goder/order/domain/order"
	"go.mongodb.org/mongo-driver/mongo"
)

// 這裡的命名有講究，先寫 orderRepository 再寫實際的實現實體名稱，這樣在其他地方要注入時
// 只要輸入 OrderRepository 就可以選擇自己想要的實體注入。
type OrderRepositoryMongo struct {
	db *mongo.Client
}

func (o *OrderRepositoryMongo) collection() *mongo.Collection {
	return o.db.Database(DBName).Collection(CollName)
}

func (o *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrderRepositoryMongo) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	//TODO implement me
	panic("implement me")
}
