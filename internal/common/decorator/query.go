package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

// QueyrHandler defines a generic type that receives a Query Q,
// and returns a result R
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

// 初始化一個 decorators，注入 log 和 metrics client
func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[H, R] {
	// 因為 queryLoggingDecorator 也實現了 QueryHandler，所以可以回傳 queryLoggingDecorator
	// 通常返回的優先級一定是 log 優先
	return queryLoggingDecorator[H, R]{
		logger: logger,
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
