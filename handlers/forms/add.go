package forms

import (
	"manage/model"

	"github.com/gin-gonic/gin"
)

//AddForm user add form
func AddForm(c *gin.Context) {
	form := model.Form{}
	c.BindJSON(&form)
	form.AuthorID = 2
	if err := model.DB.Create(&form); err != nil {
		c.JSON(500, gin.H{"msg": "创建失败"})
	} else {
		c.JSON(200, gin.H{"msg": "创建订单成功"})
	}
}
