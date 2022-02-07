package request

import "net/http"

type Form struct {}

func (F *Form) Parse(r *http.Request)map[string]string {
	params := map[string]string{}

	r.ParseForm()
   	if len(r.Form) > 0 {
      for k,v := range r.Form {
         params[k] = v[0]
      }
    }

	return params
}
