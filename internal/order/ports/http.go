package ports

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lingjun0314/goder/common"
	client "github.com/lingjun0314/goder/common/client/order"
	"github.com/lingjun0314/goder/order/app"
	"github.com/lingjun0314/goder/order/app/command"
	"github.com/lingjun0314/goder/order/app/dto"
	"github.com/lingjun0314/goder/order/app/query"
	"github.com/lingjun0314/goder/order/convertor"
)

type HTTPServer struct {
	common.BaseResponse
	app *app.Application
}

func NewHTTPServer(app *app.Application) *HTTPServer {
	return &HTTPServer{
		app: app,
	}
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerId string) {
	var (
		req client.CreateOrderRequest
		err error
		res dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, &res)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}

	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWihhQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}

	res = dto.CreateOrderResponse{
		CustomerID:  req.CustomerId,
		OrderID:     r.OrderID,
		RedirectURL: fmt.Sprintf("http://localhost:8282/sucess?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerId string, orderId string) {
	var (
		err error
		res interface{}
	)
	defer func() {
		H.Response(c, err, res)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerID: customerId,
		OrderID:    orderId,
	})
	if err != nil {
		return
	}
	res = convertor.NewOrderConvertor().EntityToClient(o)
}
