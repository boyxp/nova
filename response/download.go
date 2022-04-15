package response

import "os"
import "fmt"
import "path"
import "net/http"
import "net/url"
import "github.com/boyxp/nova/register"

type Download struct {}

func (D *Download) Render(result interface{}) {
	local, ok := result.(string)
	if !ok {
    	D.Error("文件路径必须为string", 1001)
    	return
	}

    info, err := os.Stat(local)
    if err != nil && os.IsNotExist(err) {
        D.Error("文件不存在", 1002)
        return
    }

	filesize := info.Size()
    filename := url.QueryEscape(path.Base(local))

	w := register.GetResponseWriter()
	r := register.GetRequest()
	w.Header().Add("Content-Type", "application/octet-stream")
    w.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
    w.Header().Add("Content-Length", string(filesize))
    w.Header().Add("Cache-Control", "max-age=0")
	w.Header().Add("Cache-Control", "max-age=1")
	w.Header().Add("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	w.Header().Add("Cache-Control", "cache, must-revalidate")
	w.Header().Add("Pragma", "public")
	http.ServeFile(w, r, local)
}

func (D *Download) Error(message string, code int64) {
	w := register.GetResponseWriter()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "%d:%v", code, message)
}
