package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// 日誌裝飾器也要實現 QueryHandler 接口，用來增強 QueryHandler 的功能
type queryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing query")

	//	用 defer 就可以拿到 return 回來的結果
	defer func() {
		if err == nil {
			logger.Info("Query execute successfully")
		} else {
			logger.Error("Failed to execute query:", err)
		}
	}()

	//	調用傳過來的 handler 的 Handle 方法
	return q.base.Handle(ctx, cmd)
}

func generateActionName(cmd any) string {
	//	傳過來的 cmd 可能是 query.XXX，那麼我只需要取得 XXX 即可
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
