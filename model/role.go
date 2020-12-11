package model

import (
	"gorm.io/gorm"
)

//Role model of role
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(20)" json:"role_name"`
	Description string `gorm:"type:varchar(100)" json:"description"`
}

//RolePermission Connect role and permission
type RolePermission struct {
	gorm.Model
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
}

//IsRoleExist used to detect role is whether in database
func (r *Role) IsRoleExist() (bool, error) {
	err := DB.Where("name = ?", r.Name).First(r).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if r.ID > 0 {
		return true, nil
	}
	return false, nil
}

//Save save role struct to database
func (r *Role) Save(permissions []string) error {
	tx := DB.Begin()

	if err := tx.Create(r).Error; err != nil {
		tx.Rollback()
		return err
	}

	l := len(permissions)
	rolepermissions := make([]RolePermission, l)

	for i, v := range permissions {
		permission := Permission{}
		p := v
		if err := tx.Where("name = ?", p).First(&permission).Error; err != nil {
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

//GetPermissions get role's permissions
func (r *Role) GetPermissions() ([]Permission, error) {
	var returnerr error
	pr := []RolePermission{}
	if err := DB.Where("role_id = ?", r.ID).Find(&pr).Error; err != nil {
		returnerr = err
	}

	l := len(pr)
	p := make([]Permission, l)

	for i := 0; i < l; i++ {
		if err := DB.Where("id = ?", pr[i].PermissionID).First(&p[i]).Error; err != nil {
			returnerr = err
		}
	}
	return p, returnerr
}

//HasPermisson verify whether role has permission
func (r *Role) HasPermisson(permission string) (bool, error) {
	permisions, err := r.GetPermissions()
	var returnerr error
	if err != nil {
		returnerr = err
	}
	hasP := false
	l := len(permisions)
	for i := 0; i < l; i++ {
		if permisions[i].URL == permission {
			hasP = true
			break
		}
	}
	return hasP, returnerr
}
