package response

import "net/http"

type Interface interface {
	Render(w http.ResponseWriter)
}
