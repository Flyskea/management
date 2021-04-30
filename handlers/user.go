package handlers

import (
	"manage/dto"
	"manage/middlewares"
	"manage/model"
	"manage/utils"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AddUser add a user with role
func AddUser(c *gin.Context) {
	u := dto.UR{}
	user := model.User{}
	role := &model.Role{}

	if err := c.BindJSON(&u); err != nil {
		utils.BadRequest(c, nil, "名称或角色不能为空")
		return
	}
	bytesPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		utils.InternalError(c, nil, "密码加密失败")
		return
	}
	user.Name = u.Name
	user.Password = string(bytesPwd)
	user.WorkID = u.WorkID
	role.Name = u.Role

	exist, err := user.IsUserExist()
	if err != nil {
		utils.InternalError(c, nil, "创建失败")
		return
	}
	if exist {
		utils.BadRequest(c, nil, "该用户已经存在")
		return
	}
	exist, err = role.IsRoleExist()
	if err != nil {
		utils.InternalError(c, nil, "创建失败")
		return
	}
	if !exist {
		utils.BadRequest(c, nil, "角色不存在")
		return
	}

	if err := user.Save(role.ID); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	userDTO := dto.UserDTO{}
	userDTO.Convert(&user)
	utils.Response(c, http.StatusCreated, gin.H{"user": userDTO}, "角色添加成功")
}

// UpdateUserRole 。。
func UpdateUserRole(c *gin.Context) {
	type Role struct {
		ID uint `json:"rid" binding:"required"`
	}
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
	role := Role{}
	if err := c.BindJSON(&role); err != nil {
		utils.BadRequest(c, nil, err.Error())
	}

	role1 := model.Role{}
	if err := model.DB.Where("id = ?", role.ID).First(&role1).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "角色已经删除或没有该角色")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}

	user := model.User{}
	if err := model.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "没有该用户")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if err := user.UpdateRole(role.ID); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	userDTO := dto.UserDTO{}
	userDTO.Convert(&user)
	utils.Success(c, gin.H{"user": userDTO}, "更新用户角色成功")
}

// DeleteUser delete user
func DeleteUser(c *gin.Context) {
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
	user := model.User{}
	if err := model.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "没有该用户")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if err := user.Delete(); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, nil, "删除角色成功")
}

// UserLists get user lists
func UserLists(c *gin.Context) {
	query := model.DB
	query = query.Model(&model.User{})
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
	users := []model.User{}
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	userDTOs := []dto.UserDTO{}
	for _, v := range users {
		userDTO := dto.UserDTO{}
		userDTO.Convert(&v)
		userDTOs = append(userDTOs, userDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"users":       userDTOs,
	}
	utils.Success(c, data, "查询用户列表成功")
}

// @Summary 测试用户登录
// @Description 用户登录
// @Tags user
// @Accept json
// @Produce json
// @Param name body object true "人名"
// @Param pwd body object true "密码"
// @Success 200 {string} string "{"msg": "hello Razeen"}"
// @Failure 400 {string} string "{"msg": "who are you"}"
// @Router /login [post]
//Login log in handler
func Login(c *gin.Context) {
	// 登录信息
	type LoginInfo struct {
		Name     string `json:"name" binding:"required"`
		Password string `json:"pwd" binding:"required"`
	}
	var info LoginInfo
	var user model.User
	ur := model.UserRole{}
	session := sessions.Default(c)
	sessUID := session.Get("UserID")
	if sessUID != nil {
		utils.BadRequest(c, nil, "请勿重复登陆")
		return
	}
	if err := c.BindJSON(&info); err != nil {
		utils.BadRequest(c, nil, "json信息错误")
		return
	}
	if err := model.DB.Where("Name = ?", info.Name).First(&user).Error; err != nil {
		utils.BadRequest(c, nil, "用户名错误")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {
		utils.BadRequest(c, nil, "密码错误")
		return
	}
	middlewares.GenerateNewCsrf(c, session)
	model.DB.Where("user_id = ?", user.ID).First(&ur)

	session.Set("RoleID", ur.RoleID)
	session.Set("UserID", user.WorkID)
	if err := session.Save(); err != nil {
		utils.InternalError(c, nil, "Cookie设置错误")
	}

	utils.Success(c, nil, "登录成功")
}
