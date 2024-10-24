package ports

import (
	"context"

	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct{}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (G GRPCServer) CreateOrder(context.Context, *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {

}
func (G GRPCServer) GetOrder(context.Context, *orderpb.GetOrderRequest) (*orderpb.Order, error) {

}
func (G GRPCServer) UpdateOrder(context.Context, *orderpb.Order) (*emptypb.Empty, error) {

}
