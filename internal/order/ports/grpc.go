package ports

import (
	"context"

	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app *app.Application
}

func NewGRPCServer(app *app.Application) *GRPCServer {
	return &GRPCServer{
		app: app,
	}
}

func (G GRPCServer) CreateOrder(context.Context, *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	return nil,nil
}
func (G GRPCServer) GetOrder(context.Context, *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	return nil,nil
}
func (G GRPCServer) UpdateOrder(context.Context, *orderpb.Order) (*emptypb.Empty, error) {
	return nil,nil
}
