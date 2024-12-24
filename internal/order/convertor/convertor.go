package convertor

import (
	client "github.com/lingjun0314/goder/common/client/order"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	domain "github.com/lingjun0314/goder/order/domain/order"
	"github.com/lingjun0314/goder/order/entity"
)

type OrderConvertor struct {
}

type ItemConvertor struct {
}

type ItemWithQuantityConvertor struct{}

func (c ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c ItemWithQuantityConvertor) EntityToProto(i *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ProtoToEntity(i))
	}
	return
}

func (c ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       i.Id,
		Quantity: i.Quantity,
	}
}

func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(proto *orderpb.Order) *domain.Order {
	c.check(proto)
	return &domain.Order{
		ID:          proto.ID,
		CustomerID:  proto.CustomerID,
		Status:      proto.Status,
		Items:       NewItemConvertor().ProtosToEntities(proto.Items),
		PaymentLink: proto.PaymentLink,
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
	return &domain.Order{
		ID:          o.Id,
		CustomerID:  o.CustomerId,
		Status:      o.Status,
		Items:       NewItemConvertor().ClientsToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
	c.check(o)
	return &client.Order{
		Id:          o.ID,
		CustomerId:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToClients(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) check(o interface{}) {
	if o == nil {
		panic("connot convert nil order")
	}
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, i := range items {
		res = append(res, c.ProtoToEntity(i))
	}
	return
}

func (c *ItemConvertor) ClientsToEntities(items []client.Item) (res []*entity.Item) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c *ItemConvertor) EntitiesToClients(items []*entity.Item) (res []client.Item) {
	for _, i := range items {
		res = append(res, c.EntityToClient(i))
	}
	return
}

func (c *ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(i *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ClientToEntity(i client.Item) *entity.Item {
	return &entity.Item{
		ID:       i.Id,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceId,
	}
}

func (c *ItemConvertor) EntityToClient(i *entity.Item) client.Item {
	return client.Item{
		Id:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceId:  i.PriceID,
	}
}
