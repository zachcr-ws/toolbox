package common

type Head struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type JsonResp struct {
	Head
	Data interface{} `json:"data,omitempty"`
}

func (c *JsonResp) Err(code int, msg string) {
	c.Head.Code = code
	c.Head.Msg = msg
}

func ErrResp(code int, msg string) *Head {
	data := &Head{
		Code: code,
		Msg:  msg,
	}
	return data
}

func Resp(code int, msg string, data interface{}) *JsonResp {
	ch := Head{
		Code: code,
		Msg:  msg,
	}
	result := new(JsonResp)
	result.Head = ch
	result.Data = data
	return result
}
