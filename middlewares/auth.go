package middlewares

import (
	"manage/model"
	"manage/serializer"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var normalUserPermissions map[string]struct{} = map[string]struct{}{
	"GET:/api/v1/myorder":        {},
	"GET:/api/v1/order":          {},
	"GET:/htmlselect":            {},
	"GET:/api/v1/allpermissions": {},
	"GET:/api/v1/menus":          {},
	"POST:/login":                {},
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
		c.JSON(http.StatusUnauthorized, serializer.Response{
			Msg:    "请先登录",
			Status: serializer.ErrLoginRequired,
		})
		c.Abort()
	}
}

//PermissionAuth a middleware of verifying whether user has permission to enter this url
func PermissionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestURL := c.Request.Method + ":" + c.Request.URL.Path
		if _, ok := normalUserPermissions[requestURL]; ok {
			c.Next()
			return
		}
		session := sessions.Default(c)
		currentRoleID := session.Get("RoleID").(uint)

		currentRole := model.Role{}
		if err := model.DB.Where("id = ?", currentRoleID).First(&currentRole).Error; err != nil {
			c.JSON(http.StatusInternalServerError, serializer.DBErr(err))
			c.Abort()
			return
		}
		permissions, err := currentRole.GetPermissions()
		if err != nil {
			c.JSON(http.StatusInternalServerError, serializer.DBErr(err))
			c.Abort()
			return
		}

		h, err := currentRole.HasPermissonByURL(permissions, requestURL)
		if !h {
			c.JSON(http.StatusForbidden, serializer.Response{
				Msg:    "请先登录",
				Status: serializer.ErrPermissionDenied,
			})
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, serializer.DBErr(err))
			c.Abort()
			return
		}
	}
}
