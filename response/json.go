package response

import "net/http"
import "encoding/json"
import "fmt"
import "github.com/boyxp/nova/exception"

type Json struct {
	w http.ResponseWriter
}

func (J *Json) Render(w http.ResponseWriter, result interface{}) {

}

func (J *Json) Error(w http.ResponseWriter,message string, code int64) {
	res := &exception.Exception{message,code}
	json, err := json.Marshal(res)
	if err != nil {
        fmt.Fprintf(w, "{\"message\":\"%v\",\"code\":0}", err)
    } else {
		fmt.Fprintf(w, string(json))
	}
}
