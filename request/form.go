package request

import "net/http"

type Form struct {}

func (F *Form) Parse(r *http.Request) []string {
	return []string{"a"}
}
