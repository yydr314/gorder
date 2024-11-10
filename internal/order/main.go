package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/lingjun0314/goder/common/config"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/common/server"
	"github.com/lingjun0314/goder/order/ports"
	"github.com/lingjun0314/goder/order/service"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		service := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, service)
	})

	httpServer := ports.NewHTTPServer(application)

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, httpServer, ports.GinServerOptions{
			BaseURL:     "/api",
			Middlewares: []ports.MiddlewareFunc{},
			ErrorHandler: func(*gin.Context, error, int) {
			},
		})
	})

}
