package main

import (
	"fmt"
	_ "manage/config"
	"manage/router"

	"github.com/spf13/viper"
)

// @title Swagger Example API
// @version 1.0
// @description 四川农业大学网络维修平台
// @contact.name Flyskea
// @host 127.0.0.1
// @schemes http
// @BasePath /api/v1
func main() {
	r := router.NewRouter()
	addr := viper.GetString("gin.address")
	port := viper.GetString("gin.port")
	if err := r.Run(addr + ":" + port); err != nil {
		panic(fmt.Sprintf("gin启动失败：%s", err))
	}
}
