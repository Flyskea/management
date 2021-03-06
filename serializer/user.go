package serializer

import (
	"manage/logger"
	"manage/model"
)

// User 用户序列化器
type User struct {
	ID        uint   `json:"id" example:"1"`
	RoleID    uint   `json:"role_id" example:"1"`
	WorkID    string `json:"work_id" example:"20189999"`
	CreatedAt int64  `json:"created_at" example:"1111111111"`
	DeletedAt int64  `json:"deleted_at" example:"1111111111"`
	Name      string `json:"user_name" example:"Flyskea"`
	Phone     string `json:"phone" example:"13730822098"`
}

// UserResponse 单个用户序列化
type UserResponse struct {
	Response
	Data User `json:"data"`
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	roleID := 0
	if err := model.DB.Table("users").Select("user_roles.role_id").
		Where("users.work_id = ?", user.WorkID).
		Joins("JOIN user_roles on users.id = user_roles.user_id").
		Find(&roleID).Error; err != nil {
		logger.Error(err)
	}
	delete := user.DeletedAt.Time.Unix()
	if delete < 0 {
		delete = 0
	}
	return User{
		ID:        user.ID,
		Name:      user.Name,
		WorkID:    user.WorkID,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Unix(),
		DeletedAt: delete,
		RoleID:    uint(roleID),
	}
}

func BuildUsers(items []model.User) []User {
	var (
		users []User
	)
	for _, item := range items {
		user := BuildUser(item)
		users = append(users, user)
	}
	return users
}

// BuildUserResponse 序列化用户响应
func BuildUserResponse(user model.User, msg string) UserResponse {
	return UserResponse{
		Response: Response{
			Msg: msg,
		},
		Data: BuildUser(user),
	}
}
