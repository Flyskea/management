package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//DB used to operate database
var DB *gorm.DB

//ConnectMysql connnect Mysql database
func ConnectMysql(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %s", err)
	}
	DB = db
	_ = DB.AutoMigrate(&User{}, &Permission{}, &Role{}, &RolePermission{}, &UserRole{}, &Form{}, &Score{}, &HTMLSelect{})
	return nil
}
