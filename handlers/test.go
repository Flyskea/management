package handlers

import (
	"manage/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Getp(c *gin.Context) {
	c.JSON(http.StatusOK, model.GetAllPermissions())
}
