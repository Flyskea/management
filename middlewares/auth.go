package middlewares

import (
	"manage/model"
	"manage/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var normalUserPermissions map[string]struct{} = map[string]struct{}{
	"GET:/api/v1/myorder":    {},
	"GET:/api/v1/order":      {},
	"GET:/htmlselect":        {},
	"GET:/api/v1/allpermiss": {},
	"GET:/api/v1/menus":      {},
}

//LoginAuth a middleware of verifying whether user is logged in
func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessUID := session.Get("RoleID")
		if sessUID != nil {
			c.Next()
			return
		}
		session.Clear()
		utils.Response(c, http.StatusUnauthorized, nil, "请先登录")
		c.Abort()
	}
}

//PermissionAuth a middleware of verifying whether user has permission to enter this url
func PermissionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURL := c.Request.Method + ":" + c.Request.URL.Path
		if _, ok := normalUserPermissions[requestURL]; ok {
			c.Next()
		}
		session := sessions.Default(c)
		currentRoleID := session.Get("RoleID").(uint)

		currentRole := model.Role{}
		if err := model.DB.Where("id = ?", currentRoleID).First(&currentRole).Error; err != nil {
			utils.InternalError(c, nil, "没有该角色或服务器内部错误")
			c.Abort()
			return
		}
		permissions, err := currentRole.GetPermissions()
		if err != nil {
			utils.InternalError(c, nil, "数据库操作失败")
			c.Abort()
			return
		}

		h, err := currentRole.HasPermissonByURL(permissions, requestURL)
		if !h {
			utils.Response(c, http.StatusForbidden, nil, "没有权限")
			c.Abort()
			return
		}
		if err != nil {
			utils.InternalError(c, nil, "数据库操作失败")
			c.Abort()
			return
		}
	}
}
