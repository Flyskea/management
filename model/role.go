package model

import (
	"manage/utils"

	"gorm.io/gorm"
)

//Role model of role
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20)" json:"roleName"`
	Description string `gorm:"type:varchar(100)" json:"roleDescription"`
}

//RolePermission Connect role and permission
type RolePermission struct {
	gorm.Model
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
}

//IsRoleExist used to detect role is whether in database
func (r *Role) IsRoleExist() (bool, error) {
	err := DB.Where("id= ?", r.ID).First(r).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if r.ID > 0 {
		return true, nil
	}
	return false, nil
}

//Save save role struct to database
func (r *Role) Save(permissions []uint) error {
	tx := DB.Begin()

	if err := tx.Create(r).Error; err != nil {
		tx.Rollback()
		return err
	}

	l := len(permissions)
	rolepermissions := make([]RolePermission, l)

	for i, v := range permissions {
		permission := Permission{}
		if err := tx.Where("id = ?", v).First(&permission).Error; err != nil {
			tx.Rollback()
			return err
		}
		rolepermissions[i] = RolePermission{
			PermissionID: permission.ID,
			RoleID:       r.ID,
		}
	}
	if l != 0 {
		if err := tx.Create(&rolepermissions).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

//Update Update role struct to database
func (r *Role) Update() error {
	tx := DB.Begin()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(r).Updates(r).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Delete delete role
func (r *Role) Delete() (err error) {
	tx := DB.Begin()
	if err = tx.Error; err != nil {
		return err
	}
	if err = tx.Delete(r).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Where("role_id = ?", r.ID).Delete(&RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// DeletePermission delete role's permission
func (r *Role) DeletePermission(pid uint) (err error) {
	tx := DB.Begin()
	if err = tx.Error; err != nil {
		return err
	}
	if err = tx.Where("permissions_id = ?", pid).Delete(&RolePermission{}).Error; err != nil {
		return err
	}
	return tx.Commit().Error
}

func (r *Role) getPermissionHelper(menu int) ([]*PermissionNode, error) {
	var returnerr error
	p := []Permission{}
	db := DB.Table("permissions").Select("permissions.id, permissions.url, permissions.path, permissions.component, permissions.parent_id, permissions.is_menu").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id JOIN roles ON roles.id = role_permissions.role_id")
	if r.ID == 0 {
		if menu == 0 || menu == 1 {
			db = db.Where("permissions.is_menu = ?", menu)
		}
	} else {
		if menu == 0 || menu == 1 {
			db = db.Where("roles.id = ? and permissions.is_menu = ?", r.ID, menu)
		} else {
			db = db.Where("roles.id = ?", r.ID)
		}
	}
	returnerr = db.Find(&p).Error
	pn := []*PermissionNode{}
	for _, v := range p {
		pn = append(pn, &PermissionNode{
			ID:        v.ID,
			ParentID:  v.ParentID,
			URL:       v.URL,
			Name:      v.Name,
			Path:      v.Path,
			Component: v.Component,
			IsMenu:    v.IsMenu,
			Child:     make([]*PermissionNode, 0),
		})
	}
	pMap := map[uint]*PermissionNode{}
	length := len(pn)
	for i := 0; i < length; i++ {
		node := pn[i]
		pMap[node.ID] = node
	}
	permissions := []*PermissionNode{}
	for _, v := range pn {
		if v.ParentID == 0 {
			permissions = append(permissions, v)
		} else {
			pMap[v.ParentID].Child = append(pMap[v.ParentID].Child, v)
		}
	}
	return permissions, returnerr
}

// GetPermissions get role's permissions
func (r *Role) GetPermissions() ([]*PermissionNode, error) {
	return r.getPermissionHelper(0)
}

//GetMenus get role's menus
func (r *Role) GetMenus() ([]*PermissionNode, error) {
	return r.getPermissionHelper(1)
}

// GetPermissionsAndMenus get role's all permissions
func (r *Role) GetPermissionsAndMenus() ([]*PermissionNode, error) {
	return r.getPermissionHelper(2)
}

//HasPermissonByURL verify whether role has permission by URL
func (r *Role) HasPermissonByURL(permissions []*PermissionNode, permission string) (bool, error) {
	var returnerr error
	hasP := false
	l := len(permissions)
	for i := 0; i < l; i++ {
		if utils.KeyMatch2(permission, permissions[i].URL) {
			hasP = true
			break
		}
	}
	return hasP, returnerr
}

//HasPermissionByName verify whether role has permission by name
func (r *Role) HasPermissionByName(permissions []Permission, name string) (bool, error) {
	var returnerr error
	hasP := false
	for _, v := range permissions {
		if name == v.Name {
			hasP = true
			break
		}
	}
	return hasP, returnerr
}
