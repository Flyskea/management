package model

import (
	"time"

	"gorm.io/gorm"
)

//Form Form model
type Form struct {
	gorm.Model
	//订单生成时间
	OrderedAt time.Time `gorm:"default:null"`
	//订单审核时间
	AudittedAt time.Time `gorm:"default:null"`
	//订单接收时间
	GotAt time.Time `gorm:"default:null"`
	//订单完成时间
	DoneAt time.Time `gorm:"default:null"`
	//维修人员ID
	MaintainerID uint
	//维修人员IP
	MaintainerIP string `gorm:"type:varchar(20)"`
	//评价
	Feedback float32
	//审核人
	AuditorID uint
	//订单发布人
	AuthorID uint
	//用户地址
	Location string `gorm:"type:varchar(20)" json:"location"`
	//该订单问题类型
	ProblemType string `gorm:"type:varchar(20)" json:"problemtype"`
	//订单内容
	Text string `gorm:"type:varchar(255)" json:"text"`
	//用户联系方式
	UserTelephone string `gorm:"type:varchar(11)" json:"telephone"`
}
