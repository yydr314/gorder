package main

import (
	"context"
	"github.com/lingjun0314/goder/common/broker"
	_ "github.com/lingjun0314/goder/common/config"
	"github.com/lingjun0314/goder/common/logging"
	"github.com/lingjun0314/goder/common/server"
	"github.com/lingjun0314/goder/common/tracing"
	"github.com/lingjun0314/goder/payment/infrastructure/consumer"
	"github.com/lingjun0314/goder/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

// panic 和 Fatal 的區別： panic 還會不斷冒泡回去 main 函式，然後執行完所有的 defer 函式
// Fatal 會直接跳出，不執行接下來的所有步驟
func main() {
	serviceName := viper.GetString("payment.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverType := viper.GetString("payment.server-to-run")

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)

	defer func() {
		cleanup()
	}()

	//	註冊到 consul
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

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type")
	default:
		logrus.Panic("unsupported server type")
	}

}
