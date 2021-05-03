package handlers

import (
	"manage/middlewares"
	"manage/serializer"
	service "manage/services"
	"manage/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AddUser add a user with role
func AddUser(c *gin.Context) {
	var (
		service service.UserAddService
	)
	if err := c.BindJSON(&service); err == nil {
		if user, err := service.Register(); err == nil {
			res := serializer.BuildUserResponse(user, "添加用户成功")
			c.JSON(http.StatusOK, res)
		} else {
			c.JSON(http.StatusBadRequest, err)
		}
	} else {
		c.JSON(http.StatusBadRequest, serializer.BuildErr(err, err.Error()))
	}
}

// UpdateUserRole 。。
func UpdateUserRole(c *gin.Context) {
	var (
		id      string = c.Param("id")
		service service.UserRoleUpdateService
	)
	if err := c.BindJSON(&service); err == nil {
		service.UserID = id
		if user, rsp := service.UpdateRole(); rsp == nil {
			c.JSON(http.StatusOK, serializer.BuildUserResponse(user, "更新用户角色成功"))
		} else {
			c.JSON(http.StatusBadRequest, rsp)
		}
	} else {
		c.JSON(http.StatusBadRequest, serializer.BuildErr(err, err.Error()))
	}
}

// DeleteUser delete user
func DeleteUser(c *gin.Context) {
	var (
		service service.UserDeleteService
	)
	service.ID = c.Param("id")

	if rsp := service.Delete(); rsp != nil {
		c.JSON(http.StatusInternalServerError, rsp)
	} else {
		c.JSON(http.StatusOK, serializer.Response{
			Msg: "删除该用户成功",
		})
	}
}

// UserLists get user lists
func UserLists(c *gin.Context) {
	var (
		service service.UserListService
	)
	service.Params = c.Request.URL.Query()
	if ok, rsp := service.List(); !ok {
		c.JSON(http.StatusBadRequest, rsp)
	} else {
		c.JSON(http.StatusOK, rsp)
	}
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
	var service service.UserLoginService
	session := sessions.Default(c)
	sessUID := session.Get("UserID")
	if sessUID != nil {
		utils.BadRequest(c, nil, "请勿重复登陆")
		return
	}
	if err := c.BindJSON(&service); err != nil {
		utils.BadRequest(c, nil, "json信息错误")
		return
	}
	if user, err := service.Login(); err != nil {
		c.JSON(200, err)
		return
	} else {
		middlewares.GenerateNewCsrf(c, session)
		res := serializer.BuildUserResponse(user, "登录成功")
		session.Set("UserID", user.WorkID)
		session.Set("RoleID", res.Data.RoleID)
		session.Save()
		c.JSON(http.StatusOK, res)
	}
}
