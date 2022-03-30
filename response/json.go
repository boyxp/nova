package response

import "net/http"
import "encoding/json"
import "fmt"

type Json struct {}

func (J *Json) Render(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")

	json, err := json.Marshal(result)
	if err != nil {
        fmt.Fprintf(w, "{\"code\":-1,\"message\":\"%v\",\"response\":\"\"}", err)
        return
    }

	fmt.Fprintf(w, "{\"code\":0,\"message\":\"\",\"response\":%v}", string(json))
}

func (J *Json) Error(w http.ResponseWriter, message string, code int64) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"code\":%d,\"message\":\"%v\",\"response\":\"\"}", code, message)
}
