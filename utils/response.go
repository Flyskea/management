package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpstatus int, data gin.H, msg string) {
	c.JSON(httpstatus, gin.H{"data": data, "msg": msg})
}

func Success(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, data, msg)
}

func BadRequest(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusBadRequest, data, msg)
}

func InternalError(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusInternalServerError, data, msg)
}
