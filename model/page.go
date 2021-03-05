package model

import (
	"errors"
	"manage/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	maxPageSize = 100
)

// ErrorParse if parse url query false, get this error
var ErrorParse = errors.New("参数不正确")

// Pagination used to page data
type Pagination struct {
	// the database model you want to get
	Model interface{}
	// sql conditions
	Where string
	// sql conditions
	WhereVals []interface{}
	// sql conditions rank the result
	Order strings.Builder
	// default sql order conditions
	DefaultOrder string
	// sql conditions if true get all datas
	// otherwise if deleted_at is not null the result won't be found
	Delete bool
	// Model can be orderBy lists
	OrderByList []string
}

// Page return paged data,
// page is the current page,
// pagesize is the return size,
// model is used to query table from database,
func (p *Pagination) Page(c *gin.Context) (map[string]interface{}, error) {
	var total int64
	var data interface{}
	orderBy := c.QueryArray("orderBy")
	sortBy := c.QueryArray("sortBy")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	if p.DefaultOrder == "" {
		p.DefaultOrder = "ID ASC"
	}
	if err := p.parseOrderBy(orderBy, sortBy); err != nil {
		return nil, ErrorParse
	}
	page1, err := strconv.Atoi(page)
	if err != nil {
		return nil, ErrorParse
	}

	pagesize, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, ErrorParse
	}

	totalPage := 0
	db := DB.Model(p.Model)
	if p.Delete {
		db = db.Unscoped()
	}
	if p.Where != "" {
		db = db.Where(p.Where, p.WhereVals...)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if pagesize > maxPageSize {
		pagesize = maxPageSize
	}
	if int(total)%pagesize != 0 {
		totalPage = (int(total) / pagesize) + 1
	} else {
		totalPage = (int(total) / pagesize)
	}
	if totalPage != 0 {
		if page1 > totalPage {
			page1 %= totalPage
		}
	}
	if page1 == 0 {
		page1 = 1
	}
	offset := (page1 - 1) * pagesize
	if p.Order.Len() != 0 {
		db = db.Order(p.Order.String())
	} else {
		db = db.Order(p.DefaultOrder)
	}

	switch p.Model.(type) {
	case User:
		users := []User{}
		if err := db.Limit(pagesize).Offset(offset).Find(&users).Error; err != nil {
			return nil, err
		}
		data = users
	case Role:
		roles := []Role{}
		if err := db.Limit(pagesize).Offset(offset).Find(&roles).Error; err != nil {
			return nil, err
		}
		data = roles
	case Order:
		orders := []Order{}
		if err := db.Limit(pagesize).Offset(offset).Find(&orders).Error; err != nil {
			return nil, err
		}
		data = orders
	}

	return map[string]interface{}{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": page1,
		"PageSize":    pagesize,
		"lists":       data,
	}, nil
}

// parseOrderBy used to fill Pagination field order which is used to get the ranked result from database
func (p *Pagination) parseOrderBy(orderBy, sortBy []string) error {
	for _, v := range orderBy {
		if !utils.InArray(p.OrderByList, v) {
			return ErrorParse
		}
	}
	sort := []string{"ASC", "DESC", "asc", "desc"}
	for _, v := range sortBy {
		if !utils.InArray(sort, v) {
			return ErrorParse
		}
	}
	orderByLen := len(orderBy)
	for len(sortBy) < orderByLen {
		sortBy = append(sortBy, "ASC")
	}
	for i, v := range orderBy {
		p.Order.WriteString(v)
		p.Order.WriteString(" ")
		p.Order.WriteString(sortBy[i])
		if i != orderByLen-1 {
			p.Order.WriteString(",")
		}
	}
	return nil
}
