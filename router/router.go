package router

import (
	"manage/handlers"
	"manage/middlewares"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//NewRouter used to get new router
func NewRouter() *gin.Engine {
	r := gin.Default()
	limit := middlewares.NewLimitMap()
	r.Use(middlewares.Cors(), middlewares.Limit(&middlewares.Config{TimeLimitPerAct: 5, Per: time.Second, MaxSlack: time.Second}, limit))
	store := cookie.NewStore([]byte(viper.GetString("session.cookiesecret")))
	store.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		MaxAge:   7 * 24 * 60 * 60,
	})
	r.Use(sessions.Sessions("mysession", store))
	r.NoRoute(handlers.NotFound)
	v1 := r.Group("/api/v1", middlewares.LoginAuth(), middlewares.Csrf(nil), middlewares.PermissionAuth())
	{
		v1.GET("/role", handlers.RoleLists)
		v1.POST("/role", handlers.AddRole)
		v1.PUT("/role/:id", handlers.UpdateRole)
		v1.DELETE("/role/:id", handlers.DeleteRole)

		v1.GET("/user", handlers.UserLists)
		v1.POST("/user", handlers.AddUser)
		v1.DELETE("/user/:id", handlers.DeleteUser)

		v1.GET("/myorder", handlers.MyAddOrder)
		v1.GET("/mygotorder", handlers.MyGotOrder)
		v1.GET("/order", handlers.OrderLists)
		v1.POST("/order", handlers.AddOrder)

		v1.GET("/audit", handlers.NeedAudit)
		v1.POST("/audit/:id", handlers.AuditOrder)

		v1.GET("/accept", handlers.TakeOrderGet)
		v1.POST("accept/:id", handlers.TakeOrderPost)

		v1.POST("/finish/:id", handlers.FinishOrder)

		v1.GET("/auditcommit", handlers.NeedGrade)
		v1.POST("/auditcommit/:id", handlers.AdminGradeOrder)

		v1.GET("/menus", handlers.GetRoleMenus)
		v1.GET("/allpermissions", handlers.GetAllPermissions)
	}
	r.POST("/revoke/:id", handlers.RevokeGotOrder)
	r.GET("/htmlselect", handlers.HTMLSelect)
	r.POST("/login", handlers.Login)
	r.GET("/order/:id", handlers.GetOrderByID)
	r.GET("/order/:id/status", handlers.OrderStatus)
	r.POST("/user/:id/rights", handlers.UpdateUserRole)
	r.GET("/role/:id", handlers.GetRoleByID)
	return r
}
