package dto

type Response struct {
	Error error
	Msg   string
}

func NewErrorResponse(err error, errMsg string) Response {
	return Response{
		Error: err,
		Msg:   errMsg,
	}
}

func NewOkReponse(msg string) Response {
	return Response{
		Error: nil,
		Msg:   msg,
	}
}
