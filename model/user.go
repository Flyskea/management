package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//User User model
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20)" json:"name" binding:"required"`
	Password string `gorm:"type:varchar(100)" json:"password" binding:"required"`
	WorkID   string `gorm:"unique;type:varchar(20)" json:"wid" binding:"required"`
	Phone    string `gorm:"type:varchar(20)" json:"phone" binding:"required"`
}

const (
	PassWordCost = 10
)

//IsUserExist used to detect user is whether in database
func (u *User) IsUserExist() (bool, error) {
	err := DB.Where("work_id = ?", u.WorkID).First(u).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if u.ID > 0 {
		return true, nil
	}
	return false, nil
}

//Save save user struct to database
func (u *User) Save(roleID uint) error {
	tx := DB.Begin()
	ur := UserRole{}
	ur.RoleID = roleID

	if err := tx.Create(u).Error; err != nil {
		tx.Rollback()
		return err
	}

	ur.UserID = u.ID

	if err := tx.Create(&ur).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// UpdateRole 更新用户角色
func (u *User) UpdateRole(roleID uint) error {
	tx := DB.Begin()
	ur := UserRole{}
	ur.UserID = u.ID
	ur.RoleID = roleID

	if err := tx.Model(&ur).Where("user_id = ?", u.ID).Updates(&ur).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Delete user
func (u *User) Delete() (err error) {
	tx := DB.Begin()
	if err = tx.Error; err != nil {
		return err
	}
	RoleID := UserRole{}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit().Error // err is nil; if Commit returns error update err
		}
	}()
	err = tx.Delete(u).Error
	err = tx.Where("user_id = ?", u.ID).Delete(&RoleID).Error
	return err
}

// GetUserByName GetUserByName
func (u *User) GetUserByName() (bool, error) {
	if err := DB.Where("name = ?", u.Name).First(u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}

// GetUserByWorkID GetUserByWorkID
func (u *User) GetUserByWorkID() (bool, error) {
	if err := DB.Where("work_id = ?", u.WorkID).First(u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return false, err
		} else {
			return false, nil
		}
	}
	return true, nil
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
