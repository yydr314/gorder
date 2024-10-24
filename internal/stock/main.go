package main

import (
	"github.com/lingjun0314/goder/common/genproto/stockpb"
	"github.com/lingjun0314/goder/common/server"
	"github.com/lingjun0314/goder/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	switch serverType {
	case "grpc":
		//	實現方法的服務
		service := ports.NewGRPCServer()
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			//	註冊 gRPC 的伺服器及服務
			stockpb.RegisterStockServiceServer(server, service)
		})
	case "http":
		//	TODO
	default:
		panic("unexpected server type")
	}

}
