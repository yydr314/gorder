package main

import (
	"log"

	"github.com/lingjun0314/goder/common/config"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println(viper.Get("order"))
}
