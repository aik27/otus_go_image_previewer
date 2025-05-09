package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Code   string `json:"code,omitempty"`
}

const (
	statusOK    = "OK"
	statusError = "Error"
)

func OK() Response {
	return Response{
		Status: statusOK,
	}
}

func Error(msg string, code string) Response {
	return Response{
		Status: statusError,
		Error:  msg,
		Code:   code,
	}
}
