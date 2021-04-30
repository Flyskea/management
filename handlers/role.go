package handlers

import (
	"manage/dto"
	"manage/model"
	"manage/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRoleByID get detail infomation of role
func GetRoleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil || id1 < 0 {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	role := model.Role{}
	if err := model.DB.Unscoped().Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "角色不存在")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	permissions, err := role.GetPermissionsAndMenus()
	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, gin.H{"role": role,
		"permissions": permissions}, "获取成功")
}

// RoleLists get role lists
func RoleLists(c *gin.Context) {
	query := model.DB
	query = query.Model(&model.Role{})
	query = query.Unscoped()
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	roles := []model.Role{}
	if err := query.Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}

	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"roles":       roles,
	}
	utils.Success(c, data, "查询用户列表成功")
}

//AddRole add role
func AddRole(c *gin.Context) {
	rp := dto.RP{}
	if err := c.BindJSON(&rp); err != nil {
		utils.BadRequest(c, nil, err.Error())
		return
	}

	role := &model.Role{}
	role.Name = rp.Name
	exist, err := role.IsRoleExist()

	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if exist {
		utils.BadRequest(c, nil, "角色已经存在")
		return
	}

	role.Description = rp.Description

	if err := role.Save(rp.Permission); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "权限不存在")
		} else {
			utils.InternalError(c, nil, "数据库操作失败")
		}
		return
	}

	utils.Success(c, nil, "创建成功")
}

// AddRolePermissions add role's permissions
func AddRolePermissions(c *gin.Context) {
	type Permissions struct {
		PermissionIDs []uint `json:"pids"`
	}
	p := Permissions{}
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil || id1 < 0 {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	role := model.Role{}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "角色不存在或已经删除")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if err := c.BindJSON(&p); err != nil {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	if err := role.Save(p.PermissionIDs); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "权限不存在")
		} else {
			utils.InternalError(c, nil, "数据库操作失败")
		}
		return
	}
	permissions, err := role.GetPermissionsAndMenus()
	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, gin.H{"permissions": permissions}, "更新成功")
}

// UpdateRole update role
func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	role := model.Role{}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "角色不存在")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	c.BindJSON(&role)
	if err := role.Update(); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, gin.H{"role": role}, "更新角色信息成功")
}

// DeleteRolePermission delete role's permission by id
func DeleteRolePermission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil || id1 < 0 {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	pid := c.Param("pid")
	if pid == "" {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	pidInt, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	role := model.Role{}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "角色不存在或已经删除")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}

	if err := role.DeletePermission(uint(pidInt)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "权限不存在")
		} else {
			utils.InternalError(c, nil, "数据库操作失败")
		}
		return
	}
	permissions, err := role.GetPermissionsAndMenus()
	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, gin.H{"permissions": permissions}, "更新成功")
}

// DeleteRole delete role
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil || id1 < 0 {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	role := model.Role{}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "角色已经删除")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if err := role.Delete(); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, nil, "删除角色成功")
}
