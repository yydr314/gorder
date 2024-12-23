package ports

import (
	"fmt"
	"github.com/lingjun0314/goder/common/tracing"
	"github.com/lingjun0314/goder/order/convertor"
	"net/http"

	"github.com/gin-gonic/gin"
	client "github.com/lingjun0314/goder/common/client/order"
	"github.com/lingjun0314/goder/order/app"
	"github.com/lingjun0314/goder/order/app/command"
	"github.com/lingjun0314/goder/order/app/query"
)

type HTTPServer struct {
	app *app.Application
}

func NewHTTPServer(app *app.Application) *HTTPServer {
	return &HTTPServer{
		app: app,
	}
}

func (H HTTPServer) PostCustomerCostumerIDOrders(c *gin.Context, costumerID string) {
	ctx, span := tracing.Start(c, "PostCustomerCostumerIDOrders")
	defer span.End()

	var req client.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if costumerID != req.CustomerID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "costumerID not match"})
		return
	}

	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWihhQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}

	traceID := tracing.TraceID(ctx)
	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"trace_id":     traceID,
		"customer_id":  req.CustomerID,
		"order_id":     r.OrderID,
		"redirect_url": fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
	})
}

func (H HTTPServer) GetCustomerCostumerIDOrdersOrderID(c *gin.Context, costumerID string, orderID string) {
	ctx, span := tracing.Start(c, "GetCustomerCostumerIDOrdersOrderID")
	defer span.End()
	o, err := H.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: costumerID,
		OrderID:    orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
	}
	traceID := tracing.TraceID(ctx)
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"trace_id": traceID,
		"data": gin.H{
			"Order": o,
		},
	})
}
