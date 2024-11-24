package client

import (
	"context"
	"github.com/lingjun0314/goder/common/discovery"
	"github.com/lingjun0314/goder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, nil, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for stock grpc")
	}
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, nil, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}
