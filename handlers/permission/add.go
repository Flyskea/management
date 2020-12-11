package permission

import (
	"manage/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AddPermiss add a permission
func AddPermiss(c *gin.Context) {
	p := model.Permission{}
	c.BindJSON(&p)
	if p.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "名称不能为空"})
		return
	}
	if err := model.DB.Create(&p).Error; err != nil {
		c.JSON(500, gin.H{"msg": "服务器错误"})
		return
	}
	c.JSON(200, gin.H{"msg": "权限添加成功"})
}
