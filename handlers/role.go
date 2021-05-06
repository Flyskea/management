package handlers

import (
	"manage/serializer"
	service "manage/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 单个角色详情
// @Description 单个角色详情
// @Tags role
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Success 200 {object} serializer.Response{data=serializer.Role}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role/:id [get]
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

// @Summary 角色列表
// @Description 角色列表
// @Tags role
// @Accept json
// @Produce json
// @Param page query int false "当前页数"
// @Param size query int false "每页数量"
// @Param sort query int false "排序"
// @Success 200 {object} serializer.DataList{items=[]serializer.Role}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role [get]
// RoleLists get role lists
func RoleLists(c *gin.Context) {
	var (
		service service.RoleListService
	)
	service.Params = c.Request.URL.Query()
	SendJSON(c, service.List())
}

// @Summary 增加角色
// @Description 增加角色
// @Tags role
// @Accept json
// @Produce json
// @Param addParams body service.RoleAddService true "增加角色参数"
// @Success 200 {object} serializer.Response{data=serializer.Role}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role [post]]
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

// @Summary 增加角色的权限
// @Description 增加角色的权限
// @Tags role
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Param pids body []int false "权限ID"
// @Success 200 {object} serializer.Response{data=[]serializer.Role}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role/:id/rights [post]
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

// @Summary 更新角色
// @Description 更新角色
// @Tags role
// @Accept json
// @Produce json
// @Param updateParams body service.RoleUpdateService true "更新角色描述参数"
// @Param id query int true "角色ID"
// @Success 200 {object} serializer.Response{data=[]serializer.Role}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role/:id [post]
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

// @Summary 删除角色特定权限
// @Description 删除角色特定权限
// @Tags role
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Param pid query int true "权限ID"
// @Success 200 {object} serializer.Response{data=[]serializer.Role}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role/:id/rights/:pid/delete [post]
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

// @Summary 删除角色
// @Description 删除角色
// @Tags role
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Success 200 {object} serializer.Response "{"msg", "删除角色成功"}"
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /role/:id [post]
// DeleteRole delete role
func DeleteRole(c *gin.Context) {
	var (
		service service.RoleDeleteService
	)
	service.ID = c.Param("id")
	SendJSON(c, service.Delete())
}
