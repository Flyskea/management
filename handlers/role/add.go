package role

import (
	"manage/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//RP used to bind json with role and permission
type RP struct {
	Name        string   `json:"role_name"`
	Description string   `json:"description"`
	Permission  []string `json:"permission"`
}

//AddRole add role
func AddRole(c *gin.Context) {
	rp := RP{}
	c.BindJSON(&rp)

	if rp.Name == "" || rp.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "名称或描述不能为空"})
		return
	}

	role := &model.Role{}
	role.Name = rp.Name
	exist, err := role.IsRoleExist()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		return
	}
	if exist {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "该角色已经存在"})
		return
	}

	role.Description = rp.Description

	if err := role.Save(rp.Permission); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "该权限不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "创建成功"})
}
