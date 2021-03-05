package dto

//RP used to bind json with role and permission
type RP struct {
	Name        string `json:"role_name" binding:"required,min=1,max=20"`
	Description string `json:"description" binding:"required,min=1,max=100"`
	Permission  []uint `json:"pids"`
}
