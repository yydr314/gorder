package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// 實際業務場景 client 可能會結合一些第三方套件來記錄
// 這裡的示例代碼將資料儲存在內存中
type MetricsClient interface {
	Inc(key string, value int)
}

// 這裡是記錄查詢時間的方法
type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	//	用 defer 就可以拿到 return 回來的結果
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("query.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("query.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("query.%s.failure", actionName), 1)
		}
	}()

	//	調用傳過來的 handler 的 Handle 方法
	return q.base.Handle(ctx, cmd)
}

type commandMetricsDecorator[C, R any] struct {
	base   CommandHandler[C, R]
	client MetricsClient
}

func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	//	用 defer 就可以拿到 return 回來的結果
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("command.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("command.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("command.%s.failure", actionName), 1)
		}
	}()

	//	調用傳過來的 handler 的 Handle 方法
	return q.base.Handle(ctx, cmd)
}
