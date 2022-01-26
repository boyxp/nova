package request

import "net/http"

type Interface interface {
	Parse(r *http.Request)[]string
}
