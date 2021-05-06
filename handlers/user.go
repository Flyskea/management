package handlers

import (
	"manage/middlewares"
	"manage/serializer"
	service "manage/services"
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
		if user, rsp := service.Register(); rsp == nil {
			c.JSON(http.StatusOK, serializer.BuildUserResponse(user, "添加用户成功"))
		} else {
			SendJSON(c, rsp)
		}
	} else {
		c.JSON(http.StatusBadRequest, serializer.ParamsErr(err))
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
			SendJSON(c, rsp)
		}
	} else {
		c.JSON(http.StatusBadRequest, serializer.ParamsErr(err))
	}
}

// DeleteUser delete user
func DeleteUser(c *gin.Context) {
	var (
		service service.UserDeleteService
	)
	service.ID = c.Param("id")

	if rsp := service.Delete(); rsp != nil {
		SendJSON(c, rsp)
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
	SendJSON(c, service.List())
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
		c.JSON(http.StatusBadRequest, serializer.Response{
			Status: serializer.ErrUserInfo,
			Msg:    "请勿重复登录",
		})
		return
	}
	if err := c.BindJSON(&service); err == nil {
		if user, rsp := service.Login(); rsp != nil {
			SendJSON(c, rsp)
		} else {
			middlewares.GenerateNewCsrf(c, session)
			res := serializer.BuildUserResponse(user, "登录成功")
			session.Set("UserID", user.WorkID)
			session.Set("RoleID", res.Data.RoleID)
			session.Save()
			c.JSON(http.StatusOK, res)
		}
	} else {
		c.JSON(http.StatusBadRequest, serializer.ParamsErr(err))
	}
}
