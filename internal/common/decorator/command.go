package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

// 初始化一個 decorators，注入 log 和 metrics client
func ApplyCommandDecorators[C, R any](handler CommandHandler[C, R], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C, R] {
	// 因為 queryLoggingDecorator 也實現了 QueryHandler，所以可以回傳 queryLoggingDecorator
	// 通常返回的優先級一定是 log 優先
	return queryLoggingDecorator[C, R]{
		logger: logger,
		base: queryMetricsDecorator[C, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
