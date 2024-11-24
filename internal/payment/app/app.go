package app

import "github.com/lingjun0314/goder/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}
