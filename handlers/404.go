package handlers

import (
	"manage/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFound 404返回
func NotFound(c *gin.Context) {
	utils.Response(c, http.StatusNotFound, nil, "页面未找到")
}
