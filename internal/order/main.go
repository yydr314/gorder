package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/lingjun0314/goder/common/config"
	"github.com/lingjun0314/goder/common/genproto/orderpb"
	"github.com/lingjun0314/goder/common/server"
	"github.com/lingjun0314/goder/order/ports"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		service := ports.NewGRPCServer()
		orderpb.RegisterOrderServiceServer(server, service)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
			BaseURL:     "/api",
			Middlewares: []ports.MiddlewareFunc{},
			ErrorHandler: func(*gin.Context, error, int) {
			},
		})
	})

}
