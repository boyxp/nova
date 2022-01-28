package response

import "net/http"

type Interface interface {
	Render(w http.ResponseWriter,result interface{})
	Error(w http.ResponseWriter,message string, code int64)
}
