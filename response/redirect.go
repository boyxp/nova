package response

import "fmt"
import "net/http"
import "github.com/boyxp/nova/register"

type Redirect struct {}

func (R *Redirect) Render(result interface{}) {
	url, ok := result.(string)
	if !ok {
    	R.Error("定向目标必须为string", 1001)
    	return
	}

	w := register.GetResponseWriter()
	r := register.GetRequest()

	http.Redirect(w, r, url, 302)
}

func (R *Redirect) Error(message string, code int64) {
	w := register.GetResponseWriter()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%d:%v", code, message)
}
