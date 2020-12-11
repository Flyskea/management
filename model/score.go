package model

import (
	"gorm.io/gorm"
)

//Score score model used to recode maintainer score
type Score struct {
	gorm.Model
	UserID   uint
	User     User
	FormNum  uint
	AvgPoint float32
}
