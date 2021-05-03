package service

import (
	"manage/model"
	"manage/serializer"
	"net/url"
	"strconv"

	"gorm.io/gorm"
)

// UserLoginService 管理用户登录的服务
type UserLoginService struct {
	UserName string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserAddService 管理增加用户的服务
type UserAddService struct {
	UserName string `json:"name" binding:"required"`
	WorkID   string `json:"wid" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// UserDeleteService 管理删除用户的服务
type UserDeleteService struct {
	ID string
}

// UserRoleUpdateService 管理更改用户角色的服务
type UserRoleUpdateService struct {
	UserID string
	RoleID uint `json:"rid" binding:"required"`
}

// UserListService 管理列出用户列表的服务
type UserListService struct {
	Params url.Values
}

// Login 用户登录函数
func (service *UserLoginService) Login() (model.User, *serializer.Response) {
	var user model.User
	user.Name = service.UserName
	if err := user.GetUserByName(); err != nil {
		return user, &serializer.Response{
			Msg: "账号或密码错误",
		}
	}

	if !user.CheckPassword(service.Password) {
		return user, &serializer.Response{
			Msg: "账号或密码错误",
		}
	}
	return user, nil
}

// Valid 验证表单
func (service *UserAddService) Valid() (model.User, model.Role, *serializer.Response) {
	var (
		user model.User
		role model.Role
	)
	user.WorkID = service.WorkID
	user.Name = service.UserName
	role.Name = service.Role
	exist, err := user.IsUserExist()
	if err != nil {
		return user, role, &serializer.Response{
			Msg: "创建失败",
		}
	}
	if exist {
		return user, role, &serializer.Response{
			Msg: "该用户已经存在",
		}
	}
	exist, err = role.IsRoleExist()
	if err != nil {
		return user, role, &serializer.Response{
			Msg: "创建失败",
		}
	}
	if !exist {
		return user, role, &serializer.Response{
			Msg: "该角色不存在",
		}
	}
	return user, role, nil
}

// Register 用户注册
func (service *UserAddService) Register() (model.User, *serializer.Response) {
	var (
		user model.User
		role model.Role
		resp *serializer.Response
	)
	// 表单验证
	if user, role, resp = service.Valid(); resp != nil {
		return user, resp
	}

	// 加密密码
	if err := user.SetPassword(service.Password); err != nil {
		return user, &serializer.Response{
			Msg: "密码加密失败",
		}
	}

	// 创建用户
	if err := user.Save(role.ID); err != nil {
		return user, &serializer.Response{
			Msg: "注册失败",
		}
	}

	return user, nil
}

func (service *UserDeleteService) Delete() *serializer.Response {
	var (
		user model.User
	)
	if service.ID == "" {
		return &serializer.Response{
			Msg: "参数错误",
		}
	}
	id, err := strconv.Atoi(service.ID)
	if err != nil || id < 0 {
		return &serializer.Response{
			Msg: "参数错误",
		}
	}
	if err := model.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &serializer.Response{
				Msg: "用户不存在",
			}
		}
		return serializer.BuildErr(err, err.Error())
	}
	if err := user.Delete(); err != nil {
		return serializer.BuildErr(err, err.Error())
	}
	return nil
}

func (service *UserRoleUpdateService) UpdateRole() (model.User, *serializer.Response) {
	var (
		user model.User
		role model.Role
	)
	if service.UserID == "" {
		return user, &serializer.Response{
			Msg: "参数错误",
		}
	}
	id, err := strconv.Atoi(service.UserID)
	if err != nil || id < 0 {
		return user, &serializer.Response{
			Msg: "参数错误",
		}
	}

	if err := model.DB.Where("id = ?", service.RoleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, &serializer.Response{
				Msg: "没有该角色",
			}
		}
		return user, serializer.BuildErr(err, err.Error())
	}
	if err := model.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, &serializer.Response{
				Msg: "没有该用户",
			}
		}
		return user, serializer.BuildErr(err, err.Error())
	}
	if err := user.UpdateRole(role.ID); err != nil {
		return user, serializer.BuildErr(err, err.Error())
	}
	return user, nil
}

func (service *UserListService) List() (bool, *serializer.Response) {
	var (
		page  uint64
		size  uint64
		total int64
		err   error
		users []model.User
		query *gorm.DB = model.DB.Unscoped()
	)

	page, size, query, err = parseLimitQueryParam(query.Model(&model.User{}), service.Params)
	if err != nil {
		return false, &serializer.Response{
			Msg:   "参数错误",
			Error: err.Error(),
		}
	}

	page, size, total, err = paginate(model.DB.Model(&model.User{}).Unscoped(), page, size)
	if err != nil {
		return false, &serializer.Response{
			Msg:   "参数错误",
			Error: err.Error(),
		}
	}

	query, err = parseOrderParams(query, service.Params)
	if err != nil {
		return false, &serializer.Response{
			Msg:   "参数错误",
			Error: err.Error(),
		}
	}
	if err = query.Find(&users).Error; err != nil {
		return false, &serializer.Response{
			Msg:   "数据库连接错误",
			Error: err.Error(),
		}
	}

	return true, serializer.BuildListResponse(serializer.BuildUsers(users), uint(total), uint(page), uint(size), "查询用户列表成功")
}
