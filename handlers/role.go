package handlers

import (
	"manage/serializer"
	service "manage/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetRoleByID get detail infomation of role
func GetRoleByID(c *gin.Context) {
	var (
		service service.RoleListService
	)
	id := c.Param("id")
	if role, rsp := service.One(id); rsp != nil {
		SendJSON(c, rsp)
	} else {
		if roleDTO, err := serializer.BuildRole(role); err != nil {
			SendJSON(c, serializer.DBErr(err))
		} else {
			c.JSON(http.StatusOK, serializer.Response{
				Data: roleDTO,
				Msg:  "查询角色详情成功",
			})
		}
	}
}

// RoleLists get role lists
func RoleLists(c *gin.Context) {
	var (
		service service.RoleListService
	)
	service.Params = c.Request.URL.Query()
	SendJSON(c, service.List())
}

//AddRole add role
func AddRole(c *gin.Context) {
	var (
		service service.RoleAddService
	)
	if err := c.BindJSON(&service); err != nil {
		SendJSON(c, serializer.BuildErr(err, err.Error(), serializer.ErrParams))
	} else {
		SendJSON(c, service.Add())
	}
}

// AddRolePermissions add role's permissions
func AddRolePermissions(c *gin.Context) {
	var (
		service service.RolePermissionAddService
	)
	service.ID = c.Param("id")
	if err := c.BindJSON(&service); err != nil {
		SendJSON(c, serializer.BuildErr(err, err.Error(), serializer.ErrParams))
	}
	SendJSON(c, service.Add())
}

// UpdateRole update role
func UpdateRole(c *gin.Context) {
	var (
		service service.RoleUpdateService
	)
	service.ID = c.Param("id")
	if err := c.BindJSON(&service); err != nil {
		SendJSON(c, serializer.BuildErr(err, err.Error(), serializer.ErrParams))
	}
	SendJSON(c, service.Update())
}

// DeleteRolePermission delete role's permission by id
func DeleteRolePermission(c *gin.Context) {
	var (
		service service.RolePermissionDeleteService
	)
	service.ID = c.Param("id")
	service.PID = c.Param("pid")
	if err := c.BindJSON(&service); err != nil {
		SendJSON(c, serializer.BuildErr(err, err.Error(), serializer.ErrParams))
	}
	SendJSON(c, service.Delete())
}

// DeleteRole delete role
func DeleteRole(c *gin.Context) {
	var (
		service service.RoleDeleteService
	)
	service.ID = c.Param("id")
	SendJSON(c, service.Delete())
}
