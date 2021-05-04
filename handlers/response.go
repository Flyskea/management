package handlers

import (
	"manage/serializer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendJSON(c *gin.Context, rsp *serializer.Response) {
	switch rsp.Status {
	case serializer.ErrDatabase, serializer.ErrInternal:
		c.JSON(http.StatusInternalServerError, rsp)
	case serializer.ErrParams, serializer.ErrUserInfo:
		c.JSON(http.StatusBadRequest, rsp)
	case serializer.ErrPermissionDenied:
		c.JSON(http.StatusForbidden, rsp)
	case serializer.ErrLoginRequired:
		c.JSON(http.StatusUnauthorized, rsp)
	case 0:
		c.JSON(http.StatusOK, rsp)
	}
}
