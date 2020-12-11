package main

import (
	"fmt"
	_ "manage/config"
	"manage/router"

	"github.com/spf13/viper"
)

func main() {
	r := router.NewRouter()
	addr := viper.GetString("gin.address")
	port := viper.GetString("gin.port")
	if err := r.Run(addr + ":" + port); err != nil {
		panic(fmt.Sprintf("gin启动失败：%s", err))
	}
}
