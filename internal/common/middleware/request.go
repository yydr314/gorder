package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestIn(c)
		defer requestOut(c)

		c.Next()
	}
}

func requestOut(c *gin.Context) {
	response, _ := c.Get("response")
	start, _ := c.Get("request_start")
	startTime := start.(time.Time)
	logrus.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"proc_time_ms": time.Since(startTime).Milliseconds(),
		"response":     response,
	}).Info("request_out")
}

func requestIn(c *gin.Context) {
	c.Set("request_start", time.Now())
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // body 被讀出來了，這行寫回去
	var compactJson bytes.Buffer
	_ = json.Compact(&compactJson, bodyBytes)
	logrus.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"start": time.Now().Unix(),
		"args":  compactJson.String(),
		"from":  c.RemoteIP(),
		"url":   c.Request.RequestURI,
	}).Info("request_in")
}
