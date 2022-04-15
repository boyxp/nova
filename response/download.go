package response

import "os"
import "fmt"
import "net/http"
import "github.com/boyxp/nova/register"

type Download struct {}

func (D *Download) Render(result interface{}) {
	path, ok := result.(string)
	if !ok {
    	D.Error("文件路径必须为string", 1001)
    	return
	}

    _, err := os.Stat(path)
    if err != nil && os.IsNotExist(err) {
        D.Error("文件不存在", 1002)
        return
    }

	w := register.GetResponseWriter()
	r := register.GetRequest()
	w.Header().Add("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

func (D *Download) Error(message string, code int64) {
	w := register.GetResponseWriter()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%d:%v", code, message)
}
