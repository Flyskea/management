package model

import (
	"gorm.io/gorm"
)

//UserRole Connect user and role
type UserRole struct {
	gorm.Model
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}
