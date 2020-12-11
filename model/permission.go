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

type PermissionNode struct {
	ID        uint             `json:"id"`
	ParentID  uint             `json:"pid"`
	Name      string           `json:"name"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Child     []PermissionNode `json:"child"`
}

//AllPermissions used to return to frontend
type AllPermissions struct {
	ID        uint   `json:"id"`
	Name      string `json:"p_name"`
	Path      string `json:"path"`
	Component string `json:"component"`
	ParentID  uint   `json:"pid"`
}

func GetMenus(pid int) []PermissionNode {
	var p []Permission
	DB.Where("parent_id = ? and is_menu = ?", pid, 1).Find(&p)
	var tree []PermissionNode
	for _, v := range p {
		child := GetMenus(int(v.ID))
		node := PermissionNode{
			ID:       v.ID,
			Name:     v.Name,
			ParentID: v.ParentID,
		}
		node.Child = child
		tree = append(tree, node)
	}
	return tree
}

func GetPermissions() []AllPermissions {
	var permissions []*Permission
	DB.Where("is_menu = ?", 0).Find(&permissions)
	l := len(permissions)
	allp := make([]AllPermissions, l)
	for i := 0; i < l; i++ {
		allp[i].Component = permissions[i].Component
		allp[i].ParentID = permissions[i].ParentID
		allp[i].Path = permissions[i].Path
		allp[i].Name = permissions[i].Name
		allp[i].ID = permissions[i].ID
	}
	return allp
}

type Reponse struct {
	Menus      []PermissionNode `json:"menu_permissions"`
	Permission []AllPermissions `json:"permissions"`
}

func GetAllPermissions() (r Reponse) {
	r.Menus = GetMenus(0)
	r.Permission = GetPermissions()
	return
}
