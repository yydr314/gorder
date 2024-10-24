package ports

import (
	"context"

	"github.com/lingjun0314/goder/common/genproto/stockpb"
)

type GRPCServer struct{}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (G GRPCServer) GetItems(context.Context, *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {

}

func (G GRPCServer) CheckIfItemsInSrock(context.Context, *stockpb.CheckIfItemsInSrockRequest) (*stockpb.CheckIfItemsInSrockResponse, error)
