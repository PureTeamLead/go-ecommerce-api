package dto

type Response struct {
	Error    error
	Msg      string
	Attached interface{}
}

func NewErrorResponse(err error, errMsg string) Response {
	return Response{
		Error: err,
		Msg:   errMsg,
	}
}

func NewOkReponse(msg string, attach interface{}) Response {
	return Response{
		Error:    nil,
		Msg:      msg,
		Attached: attach,
	}
}
