package main

import (
	"context"
	"github.com/lingjun0314/goder/common/broker"
	"github.com/lingjun0314/goder/common/discovery"
	"github.com/lingjun0314/goder/common/logging"
	"github.com/lingjun0314/goder/common/tracing"
	"github.com/lingjun0314/goder/order/infrastructure/consumer"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	_ "github.com/lingjun0314/goder/common/config"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/common/server"
	"github.com/lingjun0314/goder/order/ports"
	"github.com/lingjun0314/goder/order/service"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = shutdown(ctx)
	}()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	//	註冊到 consul
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		service := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, service)
	})

	httpServer := ports.NewHTTPServer(application)

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, httpServer, ports.GinServerOptions{
			BaseURL:     "/api",
			Middlewares: []ports.MiddlewareFunc{},
			ErrorHandler: func(*gin.Context, error, int) {
			},
		})
	})

}
