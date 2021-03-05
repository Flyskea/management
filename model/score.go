package model

import (
	"gorm.io/gorm"
)

//Score score model used to recode maintainer score
type Score struct {
	gorm.Model
	UserID   string
	OrderNum uint
	SumScore uint
	AvgPoint float32
}

// Save 保存用户分数
func (s *Score) Save(userID string, score uint) error {
	if err := DB.Where("user_id = ?", userID).First(s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err1 := DB.Create(&Score{UserID: userID}).Error; err1 != nil {
				return err1
			}
		} else {
			return err
		}
	}
	s.OrderNum++
	s.SumScore += score
	s.AvgPoint = float32(s.SumScore) / float32(s.OrderNum)
	if err := DB.Save(s).Error; err != nil {
		return err
	}
	return nil
}
