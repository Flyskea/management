package handlers

import (
	"manage/model"
	"manage/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//GetRoleMenus get role's menus
func GetRoleMenus(c *gin.Context) {
	r := model.Role{}
	session := sessions.Default(c)
	id := session.Get("RoleID").(uint)
	r.ID = id
	data, _ := r.GetMenus()
	utils.Success(c, gin.H{"menus": data}, "菜单查询成功")
}

//GetAllPermissions get all menus and permissions
func GetAllPermissions(c *gin.Context) {
	role := model.Role{}
	menus, err := role.GetMenus()
	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	permissions, err := role.GetPermissions()
	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, gin.H{"menus": menus, "permissions": permissions}, "权限查询成功")
}
