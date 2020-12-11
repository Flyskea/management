package router

import (
	"manage/handlers"
	"manage/handlers/forms"
	"manage/handlers/permission"
	"manage/handlers/role"
	"manage/handlers/user"
	"manage/middlewares"

	"github.com/gin-gonic/gin"
)

//NewRouter used to get new router
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/htmlselect", handlers.HTMLSelect)
	r.POST("/addform", forms.AddForm)

	v1 := r.Group("", middlewares.LoginAuth(), middlewares.PermissionAuth())
	{
		v1.POST("/user/add", user.AddUser)
		v1.POST("/role/add", role.AddRole)
		v1.POST("/permission/add", permission.AddPermiss)
	}

	r.POST("/login", handlers.Login)
	r.GET("/getp", handlers.Getp)
	return r
}
