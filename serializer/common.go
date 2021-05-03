package serializer

// Response 基础序列化器
type Response struct {
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

// DataList 基础列表结构
type DataList struct {
	Items interface{} `json:"items"`
	Page  uint        `json:"page"`
	Size  uint        `json:"size"`
	Total uint        `json:"total"`
}

// BuildListResponse 列表构建器
func BuildListResponse(items interface{}, total, page, size uint, msg string) *Response {
	return &Response{
		Data: DataList{
			Items: items,
			Total: total,
			Page:  page,
			Size:  size,
		},
		Msg: msg,
	}
}

func BuildErr(err error, msg string) *Response {
	respErr := &Response{}
	respErr.Error = err.Error()
	respErr.Msg = msg
	return respErr
}
