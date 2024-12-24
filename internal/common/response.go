package common

import (
	"github.com/gin-gonic/gin"
	"github.com/lingjun0314/goder/common/tracing"
	"net/http"
)

type BaseResponse struct{}

type response struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(c *gin.Context, err error, data interface{}) {
	if err != nil {
		base.error(c, err)
	} else {
		base.success(c, data)
	}
}

func (base *BaseResponse) success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &response{
		Errno:   0,
		Message: "success",
		Data:    data,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}

func (base *BaseResponse) error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, &response{
		Errno:   2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(c.Request.Context()),
	})
}
