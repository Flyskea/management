package handlers

import (
	"manage/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

//HTMLSelect return html select source
func HTMLSelect(c *gin.Context) {
	var s model.HTMLSelect
	name := c.Query("name")
	var err error
	if name != "" {
		err = model.DB.Where("name = ? and parent_id = ?", name, 0).First(&s).Error
	} else {
		s.ID = 0
	}
	if err == nil {
		trees := model.GetTrees(s.ID)
		c.JSON(200, model.Response{Data: trees})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
	}
}
