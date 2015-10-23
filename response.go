package common

import ()

type CommonHead struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CommonJsonResponse struct {
	CommonHead
	Data interface{} `json:"data,omitempty"`
}

func (c *CommonJsonResponse) Err(code int, msg string) {
	c.CommonHead.Code = code
	c.CommonHead.Msg = msg
}
func GenErrResponse(code int, msg string) *CommonHead {
	data := &CommonHead{
		Code: code,
		Msg:  msg,
	}
	return data
}

func GenResponse(code int, msg string, data interface{}) *CommonJsonResponse {
	ch := CommonHead{
		Code: code,
		Msg:  msg,
	}
	result := new(CommonJsonResponse)
	result.CommonHead = ch
	result.Data = data
	return result

