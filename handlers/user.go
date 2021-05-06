package handlers

import (
	"manage/middlewares"
	"manage/serializer"
	service "manage/services"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// @Summary 测试增加用户
// @Description 增加用户
// @Tags user
// @Accept json
// @Produce json
// @Param loginParams body service.UserAddService true "名字和密码"
// @Success 200 {object} serializer.Response{data=serializer.User}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /user [post]
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

// @Summary 更新用户角色
// @Description 更新用户角色
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Param rid body int true "角色ID"
// @Success 200 {object} serializer.Response{data=serializer.User}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /user/:id/role [post]
// UpdateUserRole 更新用户角色
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

// @Summary 删除用户
// @Description 删除用户
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} serializer.Response {"Msg": "删除该用户成功"}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /user/:id/delete [post]
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

// @Summary 用户列表
// @Description 用户列表
// @Tags user
// @Accept json
// @Produce json
// @Param page query int false "当前页数"
// @Param size query int false "每页数量"
// @Param sort query int false "排序"
// @Success 200 {object} serializer.DataList{items=[]serializer.User}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
// @Router /user [get]
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
// @Param loginParams body service.UserLoginService true "名字和密码"
// @Success 200 {object} serializer.Response{data=serializer.User}
// @Failure 400 {object} serializer.Response
// @Failure 500 {object} serializer.Response
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
