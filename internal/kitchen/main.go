package main

import (
	"context"
	"github.com/lingjun0314/goder/common/broker"
	"github.com/lingjun0314/goder/common/client"
	"github.com/lingjun0314/goder/common/logging"
	"github.com/lingjun0314/goder/common/tracing"
	"github.com/lingjun0314/goder/kitchen/adapters"
	"github.com/lingjun0314/goder/kitchen/infrastructure/consumer"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lingjun0314/goder/common/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("kitchen.service-name")

	//	此 context 旨在檢測超時
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	orderClient, closeFn, err := client.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closeFn()

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

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	go consumer.NewConsumer(orderGRPC).Listen(ch)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logrus.Infof("received signal, exiting...")
		os.Exit(0)
	}()
	logrus.Println("to exit, press Ctrl+C")
	select {}
}
