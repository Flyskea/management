package handlers

import (
	"manage/model"
	"manage/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginInfo struct {
	Name     string `form:"name" json:"name"`
	Password string `form:"pwd" json:"pwd"`
}

func Login(c *gin.Context) {
	var info LoginInfo
	var user model.User
	ur := model.UserRole{}
	sess := utils.GlobalSessions.SessionStart(c)
	sessUID := sess.Get("userid")
	if sessUID != nil {
		c.JSON(http.StatusFound, gin.H{"msg": "登陆成功"})
		return
	}
	if err := c.Bind(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "json信息错误"})
	}
	if err := model.DB.Where("Name = ?", info.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "用户名错误"})
		return
	}
	if info.Password != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "密码错误"})
		return
	}
	model.DB.Where("user_id = ?", user.ID).First(&ur)
	if err := sess.Set("userRole", ur.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "cookie设置错误"})
	}
	c.JSON(200, gin.H{"msg": "登陆成功"})
}
