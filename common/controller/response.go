package controller

const (
	errorBadRequest      = "bad_request"
	errorSystemError     = "system_error"
	errorBadRequestBody  = "bad_request_body"
	errorBadRequestParam = "bad_request_param"
)

type ResponseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResponseData(data interface{}) ResponseData {
	return ResponseData{
		Data: data,
	}
}

func newResponseCodeError(code string, err error) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  err.Error(),
	}
}

func newResponseCodeMsg(code, msg string) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  msg,
	}
}
