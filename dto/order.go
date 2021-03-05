package dto

import "manage/model"

//OrderDTO used to display order on frontend
type OrderDTO struct {
	ID        uint   `json:"id"`
	AuditorID string `json:"auditor_id"`
	// 维修人员
	MaintainerID string `json:"maintainer_id"`
	//订单发布人
	AuthorID string `json:"author_id"`
	//用户地址
	Location string `json:"location"`
	//用户楼栋
	Building string `json:"building"`
	//详细地址
	DetailLocation string `json:"detaillocation"`
	//该订单问题类型
	ProblemType string `json:"problemtype"`
	//订单内容
	Text string `json:"text"`
	// 订单状态
	Status int
}

//Convert used to convert order to display struct
func (o *OrderDTO) Convert(f model.Order) {
	o.AuditorID = f.AuditorID
	o.AuthorID = f.AuthorID
	o.Building = f.Building
	o.DetailLocation = f.DetailLocation
	o.ID = f.ID
	o.Location = f.Location
	o.ProblemType = f.ProblemType
	o.Text = f.Text
	o.MaintainerID = f.MaintainerID
	o.Status = f.Status()
}
