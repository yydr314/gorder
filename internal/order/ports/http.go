package ports

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lingjun0314/goder/order/app"
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

}

func (H HTTPServer) GetCustomerCostumerIDOrdersOrderID(c *gin.Context, costumerID string, orderID string) {
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		CustomerID: costumerID,
		OrderID:    orderID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    o,
	})
}
