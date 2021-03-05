package model

import (
	"time"

	"gorm.io/gorm"
)

//Order order model
type Order struct {
	gorm.Model
	//订单审核时间 用户提交后管理员审核
	AudittedAt time.Time `gorm:"default:null"`
	//订单审核时间 完成之后给维修人员与协同人员赋分
	GradedAt time.Time `gorm:"default:null"`
	//订单接收时间
	GotAt time.Time `gorm:"default:null"`
	//订单完成时间
	DoneAt time.Time `gorm:"default:null"`
	//维修人员ID
	MaintainerID string `gorm:"type:varchar(20);default:null"`
	// 退回订单维修人员ID
	BackMaintainerID string `gorm:"type:varchar(20);default:null"`
	//维修人员IP
	MaintainerIP string `gorm:"type:varchar(20)"`
	//协同人员
	HelperID string `gorm:"type:varchar(20)"`
	//维修人员得分
	MaintainerScore uint
	//协同人员得分
	HelperScore uint
	//审核人
	AuditorID string `gorm:"type:varchar(20);default:null"`
	//订单发布人
	AuthorID string `gorm:"type:varchar(20)"`
	//用户地址
	Location string `gorm:"type:varchar(100)" json:"location"`
	//用户楼栋
	Building string `gorm:"type:varchar(100)" json:"building"`
	//详细地址
	DetailLocation string `gorm:"type:varchar(100)" json:"detaillocation"`
	//该订单问题类型
	ProblemType string `gorm:"type:varchar(40)" json:"problemtype"`
	//实际的问题类型
	ActualType string `gorm:"type:varchar(40)" json:"actualtype"`
	//订单内容
	Text string `gorm:"type:varchar(255)" json:"text"`
	//维修人员完成时提交的信息
	DoneText string `gorm:"type:varchar(255)" json:"donetext"`
	//用户联系方式
	UserTelephone string `gorm:"type:varchar(11)" json:"telephone"`
	//维修人员回退订单信息
	BackText string `gorm:"type:varchar(100)" json:"backtext,index"`
	//审核不通过信息
	NotText string `gorm:"type:varchar(100)" json:"nottext,index"`
	//评价
	Feedback float32
}

// Status 获取订单状态
func (o *Order) Status() int {
	status := 0
	if o.BackText != "" {
		status = -1
		return status
	}
	if !o.AudittedAt.IsZero() {
		status++
	}
	if !o.GotAt.IsZero() {
		status++
	}
	if !o.DoneAt.IsZero() {
		status++
	}
	if !o.GradedAt.IsZero() {
		status++
	}
	return status
}
