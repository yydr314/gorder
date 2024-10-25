package service

import (
	"context"

	"github.com/lingjun0314/goder/stock/app"
)

func NewApplication(ctx context.Context) *app.Application {
	return &app.Application{}
}
