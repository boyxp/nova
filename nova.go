package nova

import "log"
import "os"
import "net/http"
import "runtime"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/request"
import "github.com/boyxp/nova/response"
import "github.com/boyxp/nova/exception"
import "github.com/fvbock/endless"

func Listen(port string) *App {
	return &App{port, &request.Form{}, &response.Json{}}
}

type App struct {
	Port string
	Request request.Interface
	Response response.Interface
}

func (A *App) Run() {
	err := endless.ListenAndServe(A.Port, A)
	if err != nil {
		log.Println(err)
	}
	log.Println("Server stopped")

	os.Exit(0)
}

func (A *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer Exception();

	match := router.Match(r.RequestURI)
	if match ==true {
		params := A.Request.Parse(r)
		result := router.Call(r.RequestURI, params)
		log.Println(result)
	}

	//===渲染结果
}

func (A *App) SetRequest(req request.Interface) *App {
	A.Request = req
	return A
}

func (A *App) SetResponse(res response.Interface) *App {
	A.Response = res
	return A
}

//异常捕获
func Exception() {
        if err :=recover();err !=nil {
        		exception := err.(*exception.Exception)
        		log.Println("异常代码：", exception.GetCode(), "异常内容：", exception.GetMessage())

                for i := 0; ; i++ {
                    pc, file, line, ok := runtime.Caller(i)
                    if !ok {
                        break
                    }
                   	log.Println(pc, file, line)
                }
        }
}
