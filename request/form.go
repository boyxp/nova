package request

import "net/http"

type Form struct {}

func (F *Form) Parse(r *http.Request) []string {
	params := make([]string,0,1)

	r.ParseForm()
	for _,v := range r.PostForm {
		if len(v) < 1 { continue }
		params = append(params, v[0])
	}

	return params
}
