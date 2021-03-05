package handlers

import (
	"fmt"
	"html"
	"manage/dto"
	"manage/model"
	"manage/utils"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderStatus 订单状态及信息
func OrderStatus(c *gin.Context) {
	type Status struct {
		Time    string
		Context string
	}
	id := c.Param("id")
	order := model.Order{}
	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "订单不存在")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	status := []*Status{}
	if !order.GradedAt.IsZero() {
		user := model.User{WorkID: order.AuditorID}
		user.GetUserByWorkID()
		status = append(status, &Status{
			Time:    order.GradedAt.Format("2006-01-02 15:04:05"),
			Context: fmt.Sprintf("%s已评分", user.Name),
		})
	}
	if !order.DoneAt.IsZero() {
		user := model.User{WorkID: order.MaintainerID}
		user.GetUserByWorkID()
		status = append(status, &Status{
			Time:    order.DoneAt.Format("2006-01-02 15:04:05"),
			Context: fmt.Sprintf("%s已完成 手机:%s", user.Name, user.Phone),
		})
	}
	if !order.GotAt.IsZero() {
		user := model.User{WorkID: order.MaintainerID}
		user.GetUserByWorkID()
		status = append(status, &Status{
			Time:    order.GotAt.Format("2006-01-02 15:04:05"),
			Context: fmt.Sprintf("%s已接单 手机:%s", user.Name, user.Phone),
		})
	}
	if !order.AudittedAt.IsZero() {
		user := model.User{WorkID: order.AuditorID}
		user.GetUserByWorkID()
		status = append(status, &Status{
			Time:    order.AudittedAt.Format("2006-01-02 15:04:05"),
			Context: fmt.Sprintf("%s已审核", user.Name),
		})
	}
	status = append(status, &Status{
		Time:    order.CreatedAt.Format("2006-01-02 15:04:05"),
		Context: fmt.Sprintf("%s进行报修", order.AuthorID),
	})
	utils.Success(c, gin.H{"status": status}, "查询状态成功")
}

// GetOrderByID order detailed info
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order := model.Order{}
	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, "订单不存在")
			return
		}
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, gin.H{"order": order}, "订单查询成功")
}

// OrderLists get orders
func OrderLists(c *gin.Context) {
	query := model.DB
	query = query.Model(&model.Order{})
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	orders := []model.Order{}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	orderDTOs := []dto.OrderDTO{}
	for _, v := range orders {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(v)
		orderDTOs = append(orderDTOs, orderDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"orders":      orderDTOs,
	}
	utils.Success(c, data, "查询所有订单成功")
}

//AddOrder user add order
func AddOrder(c *gin.Context) {
	//AddOrder used to add order or aduit order
	type AddOrder struct {
		//用户地址
		Location string `json:"location" binding:"required,min=1,max=40"`
		//用户楼栋
		Building string `json:"building" binding:"required,min=1,max=40"`
		//详细地址
		DetailLocation string `json:"detaillocation" binding:"required,min=1,max=40"`
		//该订单问题类型
		ProblemType string `json:"problemtype" binding:"required,min=1,max=40"`
		//订单内容
		Text string `json:"text" binding:"required,min=1,max=255"`
		//用户联系方式
		UserTelephone string `json:"telephone" binding:"required,min=1,max=11"`
	}
	addorder := AddOrder{}
	order := model.Order{}
	if err := c.BindJSON(&addorder); err != nil {
		utils.BadRequest(c, nil, err.Error())
		return
	}
	fmt.Println(addorder.UserTelephone)
	session := sessions.Default(c)

	if !utils.Phone(addorder.UserTelephone) {
		utils.BadRequest(c, nil, "手机号不能为空或手机号不存在")
		return
	}

	order.AuthorID = session.Get("UserID").(string)
	order.Location = addorder.Location
	order.Building = addorder.Building
	order.DetailLocation = addorder.DetailLocation
	order.ProblemType = addorder.ProblemType
	order.Text = addorder.Text
	order.UserTelephone = addorder.UserTelephone

	if !model.IsSelectExist("location", order.Location) {
		utils.BadRequest(c, nil, "地址不存在")
		return
	}
	if !model.IsSelectExist(order.Location, order.Building) {
		utils.BadRequest(c, nil, "地址不存在")
		return
	}
	if !model.IsSelectExist("problemtype", order.ProblemType) {
		utils.BadRequest(c, nil, "问题类型不存在")
		return
	}
	if err := model.DB.Create(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
	} else {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(order)
		utils.Response(c, http.StatusCreated, gin.H{"order": orderDTO}, "订单创建成功")
	}
}

//AuditOrder audit order
func AuditOrder(c *gin.Context) {
	//AddOrder used to add order or aduit order
	type AddOrder struct {
		//用户地址
		Location string `json:"location" binding:"required,min=1,max=40"`
		//用户楼栋
		Building string `json:"building" binding:"required,min=1,max=40"`
		//详细地址
		DetailLocation string `json:"detaillocation" binding:"required,min=1,max=40"`
		//该订单问题类型
		ProblemType string `json:"problemtype" binding:"required,min=1,max=40"`
		//订单内容
		Text string `json:"text" binding:"required,min=1,max=255"`
		//用户联系方式
		UserTelephone string `json:"telephone" binding:"required,min=1,max=11"`
	}
	id := c.Param("id")
	session := sessions.Default(c)
	order := model.Order{}
	auditorder := AddOrder{}
	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if !order.AudittedAt.IsZero() {
		utils.BadRequest(c, nil, "订单已审核")
		return
	}
	if err := c.BindJSON(&auditorder); err != nil {
		utils.BadRequest(c, nil, err.Error())
		return
	}

	order.Location = auditorder.Location
	order.Building = auditorder.Building
	order.DetailLocation = auditorder.DetailLocation
	order.ProblemType = auditorder.ProblemType
	order.Text = auditorder.Text
	order.UserTelephone = auditorder.UserTelephone
	order.AudittedAt = time.Now()
	order.AuditorID = session.Get("UserID").(string)

	if err := model.DB.Model(&order).Updates(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, nil, "审核成功")
}

//NeedAudit return order which isn't auditted
func NeedAudit(c *gin.Context) {
	query := model.DB
	query = query.Model(&model.Order{})
	query = query.Where("auditted_at is null")
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	orders := []model.Order{}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	orderDTOs := []dto.OrderDTO{}
	for _, v := range orders {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(v)
		orderDTOs = append(orderDTOs, orderDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"orders":      orderDTOs,
	}
	utils.Success(c, data, "查询所有订单成功")
}

// AdminGradeOrder estimate order by admin
func AdminGradeOrder(c *gin.Context) {
	type AuditCommitOrder struct {
		Score     int `json:"score" binding:"required"`
		HelpScore int `json:"helpscore" binding:"required"`
	}
	id := c.Param("id")
	order := model.Order{}
	auditorder := AuditCommitOrder{}
	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	problemscore := model.ProblemTypeScore{}
	if err := c.BindJSON(&auditorder); err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}

	if !order.GradedAt.IsZero() {
		utils.BadRequest(c, nil, "订单已审核")
		return
	}
	if auditorder.Score == -1 || auditorder.HelpScore == -1 {
		if err := model.DB.Where("name = ?", order.ProblemType).First(&problemscore).Error; err != nil {
			utils.InternalError(c, nil, "数据库操作失败")
			return
		}
		auditorder.HelpScore = int(problemscore.Score)
		auditorder.Score = int(problemscore.Score)
	}

	order.HelperScore = uint(auditorder.HelpScore)
	order.MaintainerScore = uint(auditorder.Score)
	order.GradedAt = time.Now()

	if err := model.DB.Save(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}

	score := model.Score{}
	if err := score.Save(order.MaintainerID, order.MaintainerScore); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	score = model.Score{}
	if err := score.Save(order.HelperID, order.HelperScore); err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}

	utils.Success(c, nil, "审核成功")
}

// NeedGrade return order which isn't estimated
func NeedGrade(c *gin.Context) {
	query := model.DB
	query = query.Model(&model.Order{})
	query = query.Where("done_at is not null and graded_at is null")
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	orders := []model.Order{}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	orderDTOs := []dto.OrderDTO{}
	for _, v := range orders {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(v)
		orderDTOs = append(orderDTOs, orderDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"orders":      orderDTOs,
	}
	utils.Success(c, data, "查询所有订单成功")
}

// TakeOrderGet used to display auditted order
func TakeOrderGet(c *gin.Context) {
	query := model.DB
	query = query.Model(&model.Order{})
	query = query.Where("auditted_at IS NOT NULL AND (got_at IS NULL OR (got_at IS NOT NULL AND back_maintainer_id <> maintainer_id))")
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	orders := []model.Order{}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	orderDTOs := []dto.OrderDTO{}
	for _, v := range orders {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(v)
		orderDTOs = append(orderDTOs, orderDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"orders":      orderDTOs,
	}
	utils.Success(c, data, "查询所有订单成功")
}

// TakeOrderPost maintainer uses to take auditted order
func TakeOrderPost(c *gin.Context) {
	id := c.Param("id")
	order := model.Order{}
	session := sessions.Default(c)
	mid := session.Get("UserID").(string)

	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if order.AudittedAt.IsZero() {
		utils.BadRequest(c, nil, "该订单未审核")
		return
	}

	status := order.Status()
	if status != -1 {
		if status > 1 {
			utils.BadRequest(c, nil, "其他人已接单")
			return
		}
	} else {
		if order.BackMaintainerID == mid {
			utils.BadRequest(c, nil, "不能接自己退的单")
			return
		}
	}

	order.MaintainerID = mid
	order.MaintainerIP = c.ClientIP()
	order.GotAt = time.Now()
	if err := model.DB.Model(&order).Updates(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, nil, "成功接受订单")
}

//FinishOrder maintainer used to finish order
func FinishOrder(c *gin.Context) {
	type DoneOrder struct {
		HelperID   string `json:"helpid"`
		DoneText   string `json:"donetext" binding:"required,min=1,max=255"`
		ActualType string `json:"actualtype" binding:"required,min=1,max=40"`
	}
	id := c.Param("id")
	order := model.Order{}
	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		utils.InternalError(c, nil, err.Error())
		return
	}
	if order.AudittedAt.IsZero() {
		utils.BadRequest(c, nil, "该订单未审核")
		return
	}
	if order.GotAt.IsZero() {
		utils.BadRequest(c, nil, "该订单未接受")
		return
	}
	if !order.DoneAt.IsZero() {
		utils.BadRequest(c, nil, "订单已经完成过了")
		return
	}
	session := sessions.Default(c)
	userid := session.Get("UserID").(string)

	if order.MaintainerID != userid {
		utils.BadRequest(c, nil, "不能完成他人订单")
		return
	}
	doneorder := DoneOrder{}
	user := model.User{}

	if err := c.BindJSON(&doneorder); err != nil {
		utils.BadRequest(c, nil, err.Error())
		return
	}

	if !model.IsSelectExist("problemtype", doneorder.ActualType) {
		utils.BadRequest(c, nil, "问题类型不存在")
		return
	}

	user.WorkID = doneorder.HelperID
	exist, err := user.IsUserExist()

	if err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if !exist {
		utils.BadRequest(c, nil, "协同人员不存在")
		return
	}
	order.ActualType = html.EscapeString(doneorder.ActualType)
	order.DoneText = html.EscapeString(doneorder.DoneText)
	order.HelperID = doneorder.HelperID
	order.DoneAt = time.Now()

	if err := model.DB.Model(&order).Updates(&order).Error; err != nil {
		utils.InternalError(c, nil, err.Error())
		return
	}
	utils.Success(c, nil, "完成")
}

// MyAddOrder get the login user's orders which are written by him
func MyAddOrder(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("UserID").(string)
	query := model.DB
	query = query.Model(&model.Order{})
	query = query.Where("author_id = ?", userID)
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	orders := []model.Order{}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	orderDTOs := []dto.OrderDTO{}
	for _, v := range orders {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(v)
		orderDTOs = append(orderDTOs, orderDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"orders":      orderDTOs,
	}
	utils.Success(c, data, "查询所有订单成功")
}

// MyGotOrder get the login user's orders which are got by him
func MyGotOrder(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("UserID").(string)
	query := model.DB
	query = query.Model(&model.Order{})
	query = query.Where("maintainer_id = ?", userID)
	query, err := parseOrderParams(query, c.Request.URL.Query())
	if err != nil {
		utils.BadRequest(c, nil, "参数错误")
		return
	}
	offset, limit, total, totalPage, currentPage, perPage, err := paginate(c, query)
	if err != nil {
		utils.BadRequest(c, nil, "分页参数错误")
		return
	}
	orders := []model.Order{}
	if err := query.Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	orderDTOs := []dto.OrderDTO{}
	for _, v := range orders {
		orderDTO := dto.OrderDTO{}
		orderDTO.Convert(v)
		orderDTOs = append(orderDTOs, orderDTO)
	}
	data := gin.H{
		"TotalPage":   totalPage,
		"Total":       total,
		"CurrentPage": currentPage,
		"PageSize":    perPage,
		"orders":      orderDTOs,
	}
	utils.Success(c, data, "查询所有订单成功")
}

// RevokeGotOrder 退单
func RevokeGotOrder(c *gin.Context) {
	type BackText struct {
		Text string `json:"back_text" binding:"required"`
	}

	id := c.Param("id")
	session := sessions.Default(c)
	order := model.Order{}
	currentUser := session.Get("UserID").(string)
	backText := BackText{}

	if err := model.DB.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.BadRequest(c, nil, err.Error())
			return
		}
		utils.InternalError(c, nil, err.Error())
		return
	}

	if order.MaintainerID != currentUser {
		utils.Response(c, http.StatusForbidden, nil, "不能撤回他人订单")
		return
	}

	if !order.DoneAt.IsZero() {
		utils.BadRequest(c, nil, "不能撤回已经完成的订单")
		return
	}

	if err := c.BindJSON(&backText); err != nil {
		utils.BadRequest(c, nil, "撤回文本不能为空")
	}

	order.BackMaintainerID = currentUser
	order.BackText = backText.Text

	if err := model.DB.Model(&order).Updates(&order).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if err := model.DB.Exec("UPDATE orders SET maintainer_id = NULL Where id = ?", order.ID).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	if err := model.DB.Exec("UPDATE orders SET got_at = NULL Where id = ?", order.ID).Error; err != nil {
		utils.InternalError(c, nil, "数据库操作失败")
		return
	}
	utils.Success(c, nil, "退单成功")
}

// DeleteOrder delete order
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, nil, "参数设置错误")
		return
	}
	order := model.Order{}
	if err := model.DB.Where("id = ?", id).Delete(&order).Error; err != nil {
		utils.InternalError(c, nil, err.Error())
		return
	}
	utils.Success(c, nil, "删除订单成功")
}
