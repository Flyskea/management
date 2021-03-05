package model

import (
	"gorm.io/gorm"
)

//Permission permission model
type Permission struct {
	gorm.Model
	Name      string `gorm:"type:varchar(64)" json:"p_name"`
	URL       string `gorm:"type:varchar(64)" json:"url"`
	Path      string `gorm:"type:varchar(64)" json:"path"`
	Component string `gorm:"type:varchar(64)" json:"component"`
	ParentID  uint   `json:"pid"`
	IsMenu    bool   `json:"isMenu"`
}

//PermissionNode used to return json object to frontend
type PermissionNode struct {
	ID        uint              `json:"id"`
	ParentID  uint              `json:"pid"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	URL       string            `json:"url"`
	Component string            `json:"component"`
	IsMenu    bool              `json:"is_menu"`
	Child     []*PermissionNode `json:"child"`
}
