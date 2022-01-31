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
	defer Exception(A.Response,w);

	match := router.Match(r.RequestURI)
	if match != true {
		exception.New("路由地址错误", 100)
	}

	params := A.Request.Parse(r)
	result := router.Call(r.RequestURI, params)
	log.Println(result)
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
func Exception(res response.Interface,w http.ResponseWriter) {
        if err :=recover();err !=nil {
        		//断言逻辑异常直接抛出给用户
        		exception, ok := err.(*exception.Exception)
        		if ok {
        			res.Error(w,exception.GetMessage(), exception.GetCode())
        			log.Println("逻辑异常代码：", exception.GetCode(), "逻辑异常内容：", exception.GetMessage())
        			return
        		}

        		//返回用户模糊提示
        		res.Error(w,"系统异常请联系管理员", -100)

        		//写入精确异常日志
        		log.Println("系统异常代码：-100","系统异常内容：", err)

                for i := 0; ; i++ {
                    pc, file, line, ok := runtime.Caller(i)
                    if !ok {
                        break
                    }
                   	log.Println(pc, file, line)
                }
        }
}
