package ports

import (
	"github.com/gin-gonic/gin"
	"github.com/lingjun0314/goder/order/app"
)

type HTTPServer struct {
	app *app.Application
}

func NewHTTPServer(app *app.Application) *HTTPServer {
	return &HTTPServer{
		app: app,
	}
}

func (s HTTPServer) PostCustomerCostumerIDOrders(c *gin.Context, costumerID string) {

}

func (s HTTPServer) GetCustomerCostumerIDOrdersOrderID(c *gin.Context, costumerID string, orderID string) {

}
