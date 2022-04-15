package response

type Interface interface {
	Render(result interface{})
	Error(message string, code int64)
}
