package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const defaultPerPage = 50
const maxPerPage = 100

func calculateTotalPages(perPage, total uint64) uint64 {
	pages := total / perPage
	if total%perPage > 0 {
		return pages + 1
	}
	return pages
}

func paginate(c *gin.Context, query *gorm.DB) (offset int, limit int, total int64, totalPage uint64, page uint64, perPage uint64, err error) {
	queryPage := c.Query("page")
	queryPerPage := c.Query("per_page")
	page = 1
	perPage = defaultPerPage
	if queryPage != "" {
		page, err = strconv.ParseUint(queryPage, 10, 64)
		if err != nil {
			return
		}
	}
	if queryPerPage != "" {
		perPage, err = strconv.ParseUint(queryPerPage, 10, 64)
		if err != nil {
			return
		}
	}

	if perPage > maxPerPage {
		perPage = maxPerPage
	}

	if result := query.Count(&total); result.Error != nil {
		err = result.Error
		return
	}

	offset = int((page - 1) * perPage)
	limit = int(perPage)
	totalPage = calculateTotalPages(perPage, uint64(total))
	if page > totalPage {
		page = totalPage
	}
	return
}
