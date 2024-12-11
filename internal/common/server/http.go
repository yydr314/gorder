package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 封裝一個方法，開啟 http 服務器
func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		panic("empty http address")
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	//	開啟一個 gin 的 router (和以前我使用的 gin.Default()意思一樣)
	apiRoutper := gin.New()
	//	透過 wrapper 函式來對這個 router 進行一些配置，這樣就不需要每個地方都寫重複的代碼
	wrapper(apiRoutper)
	//	建立一個路由組
	apiRoutper.Group("/api")
	if err := apiRoutper.Run(addr); err != nil {
		panic(err)
	}
}
