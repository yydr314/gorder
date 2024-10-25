package ports

import (
	"context"

	"github.com/lingjun0314/goder/common/genproto/stockpb"
	"github.com/lingjun0314/goder/stock/app"
)

// 注入 app
type GRPCServer struct {
	app *app.Application
}

func NewGRPCServer(app *app.Application) *GRPCServer {
	return &GRPCServer{
		app: app,
	}
}

func (G GRPCServer) GetItems(context.Context, *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	return nil, nil
}

func (G GRPCServer) CheckIfItemsInSrock(context.Context, *stockpb.CheckIfItemsInSrockRequest) (*stockpb.CheckIfItemsInSrockResponse, error) {
	return nil, nil
}
