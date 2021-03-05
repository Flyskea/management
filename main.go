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
	if err := r.RunTLS(addr+":"+port, "./localhost.crt", "./localhost.key"); err != nil {
		panic(fmt.Sprintf("gin启动失败：%s", err))
	}
}
