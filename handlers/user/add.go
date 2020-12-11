package user

import (
	"manage/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//UR used to bind json data to add user and it's role
type UR struct {
	Name     string `gorm:"type:varchar(20)" json:"name"`
	WorkID   string `json:"wid"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

//AddUser add a user with role
func AddUser(c *gin.Context) {
	u := UR{}
	user := model.User{}
	role := &model.Role{}

	c.BindJSON(&u)

	if u.Name == "" || u.Password == "" || u.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "名称或角色不能为空"})
		return
	}

	user.Name = u.Name
	user.Password = u.Password
	user.WorkID = u.WorkID
	role.Name = u.Role

	exist, err := user.IsUserExist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "该用户已经存在"})
		return
	}
	exist, err = role.IsRoleExist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		return
	}
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "没有该角色"})
		return
	}

	if err := user.Save(role.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
}
