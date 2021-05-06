package serializer

import (
	"manage/model"
)

type Role struct {
	ID          uint          `json:"rid"`
	DeletedAt   int64         `json:"deleted_at"`
	Name        string        `json:"roleName"`
	Description string        `json:"roleDescription"`
	Permissions []*Permission `json:"permissions"`
}

func BuildRole(role model.Role) (r Role, err error) {
	delete := role.DeletedAt.Time.Unix()
	if delete < 0 {
		delete = 0
	}
	permissions, err := BuildPermissions(role.ID)
	return Role{
		ID:          role.ID,
		DeletedAt:   delete,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
	}, err
}

func BuildRoles(items []model.Role) ([]Role, error) {
	var (
		roles []Role
	)
	for _, item := range items {
		if role, err := BuildRole(item); err != nil {
			return nil, err
		} else {
			roles = append(roles, role)
		}
	}
	return roles, nil
}
