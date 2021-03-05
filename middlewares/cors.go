package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Cors used to cross site
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "http://192.168.1.2:8080")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Set-Cookie,Content-Type,Access-Control-Allow-Headers,Content-Length,Accept,Authorization,X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
