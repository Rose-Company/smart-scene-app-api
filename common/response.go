package common

type LOGIC_CODE uint

const (
	REQUEST_SUCCESS LOGIC_CODE = 0
	REQUEST_FAILED  LOGIC_CODE = 1
)

type Response struct {
	Code        LOGIC_CODE  `json:"code"`
	ErrorCode   string      `json:"error_code,omitempty"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
	Paging      *paging     `json:"paging,omitempty"`
	ErrorDetail string      `json:"error_detail,omitempty"`
}
type ResponseCustomize struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func ResponseCustom(code LOGIC_CODE, data interface{}, message string) Response {
	return Response{Data: data, Code: code, Message: message}
}

func ResponseSuccess(code LOGIC_CODE, data interface{}, message string) Response {
	return ResponseCustom(code, data, message)
}

func ResponseOk(data interface{}) Response {
	return ResponseCustom(0, data, "")
}

func BaseResponse(code LOGIC_CODE, message string, errorDetail string, data interface{}) Response {
	return Response{Code: code, Message: message, ErrorDetail: errorDetail, Data: data}
}

func BaseResponseMess(code int, message string, data interface{}) ResponseCustomize {
	return ResponseCustomize{Code: code, Message: message, Data: data}
}

func ResponseUnAuthorized(message string) Response {
	return Response{Code: REQUEST_FAILED, Message: message}
}

type paging struct {
	TotalCount uint `json:"total_count"`
	Limit      int  `json:"limit"`
	Page       int  `json:"page"`
}

func (t *Response) AppendPaging(totalCount uint, limit int, page int) {
	t.Paging = &paging{
		TotalCount: totalCount,
		Limit:      limit,
		Page:       page,
	}
}

func (t *Response) SetErrorCode(code string) {
	t.ErrorCode = code
}
