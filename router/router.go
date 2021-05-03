package router

import (
	_ "manage/docs"
	"manage/handlers"
	"manage/middlewares"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

//NewRouter used to get new router
func NewRouter() *gin.Engine {
	r := gin.Default()
	limit := middlewares.NewLimitMap()
	r.Use(middlewares.Cors(), middlewares.Limit(&middlewares.Config{TimeLimitPerAct: 5, Per: time.Second, MaxSlack: time.Second}, limit))
	store := cookie.NewStore([]byte(viper.GetString("session.cookiesecret")))
	store.Options(sessions.Options{
		HttpOnly: false,
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/",
	})
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.NoRoute(handlers.NotFound)
	v1 := r.Group("/api/v1", middlewares.LoginAuth(), middlewares.Csrf(nil), middlewares.PermissionAuth())
	{
		v1.GET("/htmlselect", handlers.HTMLSelect)
		roles := v1.Group("/role")
		{
			roles.GET(":id", handlers.GetRoleByID)
			roles.GET("", handlers.RoleLists)
			roles.POST("", handlers.AddRole)
			roles.POST(":id/update", handlers.UpdateRole)
			roles.POST(":id/delete", handlers.DeleteRole)
			roles.POST(":id/rights", handlers.AddRolePermissions)
			roles.POST(":id/rights/:pid/delete", handlers.DeleteRolePermission)
		}
		users := v1.Group("/user")
		{
			users.POST(":id/role", handlers.UpdateUserRole)
			users.GET("", handlers.UserLists)
			users.POST("", handlers.AddUser)
			users.POST(":id/delete", handlers.DeleteUser)
		}

		v1.GET("/order/:id", handlers.GetOrderByID)
		v1.GET("/order/:id/status", handlers.OrderStatus)
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

		v1.POST("/revoke/:id", handlers.RevokeGotOrder)

		v1.GET("/menus", handlers.GetRoleMenus)
		v1.GET("/allpermissions", handlers.GetAllPermissions)
	}
	r.POST("/api/v1/login", handlers.Login)
	return r
}
