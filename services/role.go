package service

import (
	"manage/model"
	"manage/serializer"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

// RoleListService 管理列出角色列表的服务
type RoleListService struct {
	Params url.Values
}

// RoleAddService 管理增加角色的服务
type RoleAddService struct {
	Name        string `json:"roleName" binding:"required,min=1,max=20" example:"管理员"`
	Description string `json:"description" binding:"required,min=1,max=100" example:"拥有所有权限"`
	Permission  []uint `json:"pids" example:"1,2,3,4"`
}

// RolePermissionAddService 管理增加角色权限的服务
type RolePermissionAddService struct {
	ID         string
	Permission []uint `json:"pids"`
}

// RoleUpdateService 管理更改角色名字或描述的服务
type RoleUpdateService struct {
	ID          string
	Name        string `json:"roleName" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// RoleDeleteService 管理删除角色的服务
type RoleDeleteService struct {
	ID string
}

// RolePermissionDeleteService 管理删除角色权限的服务
type RolePermissionDeleteService struct {
	ID  string
	PID string
}

func (service *RoleListService) One(rid string) (model.Role, *serializer.Response) {
	var role model.Role
	if rid == "" {
		return role, &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "参数错误",
		}
	}
	id, err := strconv.Atoi(rid)
	if err != nil || id < 0 {
		return role, serializer.ParamsErr(err)
	}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return role, &serializer.Response{
				Status: serializer.ErrParams,
				Msg:    "没有该角色",
			}
		}
		return role, serializer.DBErr(err)
	}
	return role, nil
}

func (service *RoleListService) List() *serializer.Response {
	var (
		page     uint64
		size     uint64
		total    int64
		err      error
		roles    []model.Role
		roleDTOs []serializer.Role
		query    *gorm.DB = model.DB.Unscoped()
	)

	page, size, query, err = parseLimitQueryParam(query.Model(&model.Role{}), service.Params)
	if err != nil {
		return serializer.ParamsErr(err)
	}

	page, size, total, err = paginate(model.DB.Model(&model.Role{}).Unscoped(), page, size)
	if err != nil {
		return serializer.DBErr(err)
	}

	query, err = parseOrderParams(query, service.Params)
	if err != nil {
		return serializer.ParamsErr(err)
	}
	if err = query.Find(&roles).Error; err != nil {
		return serializer.DBErr(err)
	}
	if roleDTOs, err = serializer.BuildRoles(roles); err != nil {
		return serializer.DBErr(err)
	}
	return serializer.BuildListResponse(roleDTOs, uint(total), uint(page), uint(size), "查询角色列表成功")
}

func (service *RoleAddService) valid() (model.Role, *serializer.Response) {
	var (
		role model.Role
	)
	role.Name = service.Name
	exist, err := role.IsRoleExist()

	if err != nil {
		return role, serializer.DBErr(err)
	}
	if exist {
		return role, &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "该角色已经存在",
		}
	}
	role.Description = service.Description
	return role, nil
}

func (service *RoleAddService) Add() *serializer.Response {
	if role, rsp := service.valid(); rsp != nil {
		return rsp
	} else {
		if err := role.Save(service.Permission); err != nil {
			if err == gorm.ErrRecordNotFound {
				return &serializer.Response{
					Status: serializer.ErrParams,
					Msg:    "权限不存在",
				}
			} else {
				return serializer.DBErr(err)
			}
		} else {
			if roleDTO, buildError := serializer.BuildRole(role); buildError != nil {
				return serializer.DBErr(err)
			} else {
				return &serializer.Response{
					Data: roleDTO,
					Msg:  "添加角色成功",
				}
			}
		}
	}
}

func (service *RolePermissionAddService) Add() *serializer.Response {
	var role model.Role
	if service.ID == "" {
		return &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "参数错误",
		}
	}
	id, err := strconv.Atoi(service.ID)
	if err != nil || id < 0 {
		return serializer.ParamsErr(err)
	}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.BuildErr(err, "角色不存在", serializer.ErrNotFound)
		}
		return serializer.DBErr(err)
	}

	if err := role.Save(service.Permission); err != nil {
		if err == gorm.ErrRecordNotFound {
			serializer.BuildErr(err, "权限不存在", serializer.ErrNotFound)
		}
		return serializer.DBErr(err)
	}
	if roleDTO, err := serializer.BuildRole(role); err != nil {
		return serializer.DBErr(err)
	} else {
		return &serializer.Response{
			Data: roleDTO,
			Msg:  "增加角色权限成功",
		}
	}
}

func (service *RoleUpdateService) Update() *serializer.Response {
	var role model.Role
	if service.ID == "" {
		return &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "参数错误",
		}
	}
	id, err := strconv.Atoi(service.ID)
	if err != nil || id < 0 {
		return serializer.ParamsErr(err)
	}
	if err := role.Update(); err != nil {
		return serializer.DBErr(err)
	}
	if roleDTO, err := serializer.BuildRole(role); err != nil {
		return serializer.DBErr(err)
	} else {
		return &serializer.Response{
			Data: roleDTO,
			Msg:  "更新角色名字或描述成功",
		}
	}
}

func (service *RoleDeleteService) Delete() *serializer.Response {
	var role model.Role
	if service.ID == "" {
		return &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "参数错误",
		}
	}
	id, err := strconv.Atoi(service.ID)
	if err != nil || id < 0 {
		return serializer.ParamsErr(err)
	}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.BuildErr(err, "角色不存在", serializer.ErrNotFound)
		}
		return serializer.DBErr(err)
	}
	if err := role.Delete(); err != nil {
		return serializer.DBErr(err)
	}
	return &serializer.Response{
		Msg: "删除角色成功",
	}
}

func (service *RolePermissionDeleteService) Delete() *serializer.Response {
	var role model.Role
	if service.ID == "" {
		return &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "参数错误",
		}
	}
	id, err := strconv.ParseUint(service.ID, 10, 64)
	if err != nil {
		return serializer.ParamsErr(err)
	}
	if service.PID == "" {
		return &serializer.Response{
			Status: serializer.ErrParams,
			Msg:    "参数错误",
		}
	}
	pid, err := strconv.ParseUint(service.PID, 10, 64)
	if err != nil {
		return serializer.ParamsErr(err)
	}
	if err := model.DB.Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.BuildErr(err, "角色不存在", serializer.ErrNotFound)
		}
		return serializer.DBErr(err)
	}
	if err := role.DeletePermission(uint(pid)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return serializer.BuildErr(err, "权限不存在", serializer.ErrNotFound)
		}
		return serializer.DBErr(err)
	}
	if roleDTO, err := serializer.BuildRole(role); err != nil {
		return serializer.DBErr(err)
	} else {
		return &serializer.Response{
			Data: roleDTO,
			Msg:  "删除角色权限成功",
		}
	}
}
