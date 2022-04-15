package response

import "fmt"
import "encoding/json"
import "github.com/boyxp/nova/register"

type Json struct {}

func (J *Json) Render(result interface{}) {
	w := register.GetResponseWriter()
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(result)
	if err != nil {
        fmt.Fprintf(w, "{\"code\":-1,\"message\":\"%v\",\"response\":\"\"}", err)
        return
    }

	fmt.Fprintf(w, "{\"code\":0,\"message\":\"\",\"response\":%v}", string(json))
}

func (J *Json) Error(message string, code int64) {
	w := register.GetResponseWriter()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"code\":%d,\"message\":\"%v\",\"response\":\"\"}", code, message)
}
