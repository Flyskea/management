package dto

import (
	"manage/logger"
	"manage/model"
	"time"
)

//UR used to bind json data to add user and it's role
type UR struct {
	Name     string `gorm:"type:varchar(20)" json:"name" binding:"required"`
	WorkID   string `json:"wid" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// UserDTO 返回用户列表
type UserDTO struct {
	ID        uint
	CreatedAt time.Time
	Name      string
	WorkID    string
	Phone     string
	RoleName  string
}

// Convert use user data fill dto data
func (uDTO *UserDTO) Convert(u *model.User) {
	uDTO.CreatedAt = u.CreatedAt
	uDTO.ID = u.ID
	uDTO.Name = u.Name
	uDTO.Phone = u.Phone
	uDTO.WorkID = u.WorkID
	roleName := ""
	if err := model.DB.Table("users").Select("roles.name").
		Where("users.name = ?", u.Name).
		Joins("JOIN user_roles on users.id = user_roles.user_id JOIN roles ON user_roles.role_id = roles.id").
		Find(&roleName).Error; err != nil {
		logger.Error(err)
	}
	uDTO.RoleName = roleName
}
