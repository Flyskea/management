package middlewares

import (
	"manage/model"
	"manage/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//LoginAuth a middleware of verifying whether user is logged in
func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := utils.GlobalSessions.SessionStart(c)
		sessUID := sess.Get("userRole")
		if sessUID != nil {
			c.Next()
			return
		}
		utils.GlobalSessions.SessionDestroy(c)
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登陆"})
		c.Abort()
	}
}

//PermissionAuth a middleware of verifying whether user has permission to enter this url
func PermissionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := utils.GlobalSessions.SessionStart(c)
		currentroleID := sess.Get("userRole").(uint)

		currentrole := model.Role{}
		if err := model.DB.Where("id = ?", currentroleID).First(&currentrole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "验证失败"})
			c.Abort()
			return
		}

		h, err := currentrole.HasPermisson(c.Request.Method + ":" + c.Request.RequestURI)
		if !h {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "没有权限"})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "验证失败"})
			c.Abort()
			return
		}
	}
}
