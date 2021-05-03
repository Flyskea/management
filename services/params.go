package service

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type sortDirection string

const ascending sortDirection = "asc"
const descending sortDirection = "desc"

var sortFields = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func sortField(value string) string {
	return sortFields[value]
}

func parseLimitQueryParam(query *gorm.DB, params url.Values) (uint64, uint64, *gorm.DB, error) {
	var (
		page uint64
		size uint64
		err  error
	)
	if values, exists := params["page"]; exists {
		page, err = strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return 0, 0, nil, err
		}
		query = query.Offset(int(page) - 1)
	}
	if values, exists := params["size"]; exists {
		size, err = strconv.ParseUint(values[0], 10, 64)
		if err != nil {
			return 0, 0, nil, err
		}
		query = query.Limit(int(size))
	}
	return page, size, query, nil
}

func parseOrderParams(query *gorm.DB, params url.Values) (*gorm.DB, error) {
	if values, exists := params["sort"]; exists {
		for _, value := range values {
			parts := strings.Split(value, " ")
			field := sortField(parts[0])
			if field == "" {
				return nil, fmt.Errorf("bad field for sort '%v'", field)
			}
			dir := ascending
			if len(parts) == 2 {
				switch strings.ToLower(parts[1]) {
				case string(ascending):
					dir = ascending
				case string(descending):
					dir = descending
				default:
					return nil, fmt.Errorf("bad direction for sort '%v', only 'asc' and 'desc' allowed", parts[1])
				}
			}
			query = query.Order(field + " " + string(dir))
		}
	} else {
		query = query.Order("created_at desc")
	}

	return query, nil
}
