package model

import (
	"gorm.io/gorm"
)

//User User model
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20)" json:"name"`
	Password string `gorm:"type:varchar(40)" json:"password"`
	WorkID   string `json:"wid"`
}

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
