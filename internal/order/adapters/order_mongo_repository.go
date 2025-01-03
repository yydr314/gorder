package adapters

import (
	"context"
	_ "github.com/lingjun0314/goder/common/config"
	domain "github.com/lingjun0314/goder/order/domain/order"
	"github.com/lingjun0314/goder/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	dbName   = viper.GetString("mongo.databases.db-name")
	collName = viper.GetString("mongo.databases.coll-name")
)

// 這裡的命名有講究，先寫 orderRepository 再寫實際的實現實體名稱，這樣在其他地方要注入時
// 只要輸入 OrderRepository 就可以選擇自己想要的實體注入。
type OrderRepositoryMongo struct {
	db *mongo.Client
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

func (o *OrderRepositoryMongo) collection() *mongo.Collection {
	return o.db.Database(dbName).Collection(collName)
}

type orderModel struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `bson:"customer_id"`
	Status      string             `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (o *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	defer logWithTag("create ", err, created)
	write := o.marshalToModel(order)
	res, err := o.collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}

	created = order
	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (o *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	defer logWithTag("get", err, got)

	read := &orderModel{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	// condition
	cond := bson.M{"_id": mongoID}
	err = o.collection().FindOne(ctx, cond).Decode(read)
	if err != nil {
		return nil, err
	}

	if read == nil {
		return nil, domain.NotFoundError{OrderID: id}
	}

	return o.unmarshal(read), nil
}

func (o *OrderRepositoryMongo) Update(
	ctx context.Context,
	order *domain.Order,
	updateFn func(context.Context, *domain.Order) (*domain.Order, error)) (err error) {
	if order == nil {
		panic("got nil order")
	}
	defer logWithTag("update", err, nil)
	// 開啟 mongodb 的事務，在 mongodb 中需要先開啟 session 才能開啟事務
	session, err := o.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return
	}
	defer func() {
		if err != nil {
			_ = session.AbortTransaction(ctx)
			return
		}

		_ = session.CommitTransaction(ctx)
	}()

	// inside transaction
	oldOrder, err := o.Get(ctx, order.ID, order.CustomerID)
	if err != nil {
		return
	}
	updated, err := updateFn(ctx, order)
	if err != nil {
		return
	}
	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
	res, err := o.collection().UpdateOne(
		ctx,
		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}

	logWithTag("finish_update", err, res)
	return
}

func (o *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
	return &orderModel{
		MongoID:     primitive.NewObjectID(),
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (o *OrderRepositoryMongo) unmarshal(read *orderModel) *domain.Order {
	return &domain.Order{
		ID:          read.MongoID.Hex(),
		CustomerID:  read.CustomerID,
		Status:      read.Status,
		PaymentLink: read.PaymentLink,
		Items:       read.Items,
	}
}

func logWithTag(tag string, err error, result interface{}) {
	l := logrus.WithFields(logrus.Fields{
		"tag":            "order repository mongo",
		"performed_time": time.Now().Unix(),
		"error":          err,
		"result":         result,
	})
	if err != nil {
		l.Warnf("%s failed", tag)
	} else {
		l.Infof("%s success", tag)
	}
}
