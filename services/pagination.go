package service

import (
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

func paginate(query *gorm.DB, queryPage, queryPerPage uint64) (uint64, uint64, int64, error) {
	var (
		page  uint64
		size  uint64
		total int64
		err   error
	)
	if queryPerPage == 0 {
		size = defaultPerPage
	} else if queryPerPage > maxPerPage {
		size = maxPerPage
	} else {
		size = queryPerPage
	}

	if err = query.Count(&total).Error; err != nil {
		return 0, 0, 0, err
	}
	if page == 0 {
		page = 1
	}
	tmp := calculateTotalPages(size, uint64(total))
	if queryPage == 0 {
		page = 1
	} else if queryPage > tmp {
		page = tmp
	} else {
		page = queryPage
	}
	return page, size, total, nil
}
