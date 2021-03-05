package model

import (
	"gorm.io/gorm"
)

//ProblemTypeScore databasemodel used to score the order
type ProblemTypeScore struct {
	gorm.Model
	Score uint
	Name  string `gorm:"type:varchar(64)"`
}

//db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan()
