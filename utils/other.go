package utils

import "regexp"

// InArray if value in the arr, return true.Otherwise false
func InArray(arr []string, value string) bool {
	inarr := false

	for _, v := range arr {
		if v == value {
			inarr = true
			break
		}
	}

	return inarr
}

// Phone 正则检查
func Phone(phone string) bool {
	if m, _ := regexp.MatchString(`^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`, phone); !m {
		return false
	}
	return true
}
