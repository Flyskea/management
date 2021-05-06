package serializer

import "manage/model"

type Permission struct {
	ID        uint          `json:"id"`
	ParentID  uint          `json:"pid"`
	Level     int           `json:"level"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	URL       string        `json:"url"`
	Component string        `json:"component"`
	IsMenu    bool          `json:"is_menu"`
	Child     []*Permission `json:"child"`
}

func BuildPermission(permission model.Permission) *Permission {
	return &Permission{
		ID:        permission.ID,
		ParentID:  permission.ParentID,
		Level:     permission.Level,
		Name:      permission.Name,
		Path:      permission.Path,
		URL:       permission.URL,
		Component: permission.Component,
		IsMenu:    permission.IsMenu,
		Child:     make([]*Permission, 0),
	}
}

func BuildPermissions(roleID uint) ([]*Permission, error) {
	var (
		returnerr error
		p         []model.Permission
		pn        []*Permission
		pMap      map[uint]*Permission = make(map[uint]*Permission)
	)
	db := model.DB.Table("permissions").Select("permissions.name,permissions.id, permissions.url, permissions.path, permissions.component, permissions.parent_id, permissions.is_menu,permissions.level").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id JOIN roles ON roles.id = role_permissions.role_id")
	returnerr = db.Where("roles.id = ?", roleID).Find(&p).Error

	for _, v := range p {
		pn = append(pn, BuildPermission(v))
	}
	for _, v := range pn {
		pMap[v.ID] = v
	}
	permissions := []*Permission{}
	for _, v := range pn {
		if v.ParentID == 0 {
			permissions = append(permissions, v)
		} else {
			pMap[v.ParentID].Child = append(pMap[v.ParentID].Child, v)
		}
	}
	return permissions, returnerr
}
