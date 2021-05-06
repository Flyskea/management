package serializer

const (
	ErrParams           = 40000
	ErrLoginRequired    = 40001
	ErrUserInfo         = 40002
	ErrPermissionDenied = 40003
	ErrNotFound         = 40004
	ErrInternal         = 50000
	ErrDatabase         = 50001
)

func BuildErr(err error, msg string, status int) *Response {
	respErr := &Response{}
	respErr.Error = err.Error()
	respErr.Msg = msg
	respErr.Status = status
	return respErr
}

func DBErr(err error) *Response {
	return BuildErr(err, "数据库操作失败", ErrDatabase)
}

func ParamsErr(err error) *Response {
	return BuildErr(err, "参数错误", ErrParams)
}
