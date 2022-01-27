package exception

func New(message string, code int64) {
	panic(&Exception{message, code})
}

type Exception struct {
	message string
	code int64
}

func (E *Exception) GetCode() int64 {
	return E.code
}

func (E *Exception) GetMessage() string {
	return E.message
}
