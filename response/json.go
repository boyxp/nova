package response

import "net/http"
import "encoding/json"
import "fmt"
import "github.com/boyxp/nova/exception"

type Json struct {
	w http.ResponseWriter
}

func (J *Json) Render(w http.ResponseWriter, result interface{}) {
	json, err := json.Marshal(result)

	if err != nil {
        fmt.Fprintf(w, "{\"message\":\"%v\",\"code\":0}", err)
        return
    }

	fmt.Fprintf(w, string(json))
}

func (J *Json) Error(w http.ResponseWriter,message string, code int64) {
	J.Render(w, &exception.Exception{message,code})
}
