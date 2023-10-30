package exception

func New(message string, code int64) {
	panic(&Exception{message, code})
}

func Throw(message string, code int64) {
	panic(&Exception{message, code})
}

type Exception struct {
	Message string `json:"message"`
	Code int64 `json:"code"`
}

func (E *Exception) GetCode() int64 {
	return E.Code
}

func (E *Exception) GetMessage() string {
	return E.Message
}
