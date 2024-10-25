package main

import (
	"context"

	"github.com/lingjun0314/goder/common/config"
	"github.com/lingjun0314/goder/common/genproto/stockpb"
	"github.com/lingjun0314/goder/common/server"
	"github.com/lingjun0314/goder/stock/ports"
	"github.com/lingjun0314/goder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	//	此 context 旨在檢測超時
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)
	switch serverType {
	case "grpc":
		//	實現方法的服務
		service := ports.NewGRPCServer(application)
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			//	註冊 gRPC 的伺服器及服務
			stockpb.RegisterStockServiceServer(server, service)
		})
	case "http":
		//	TODO
	default:
		//panic("unexpected server type")
	}

}
