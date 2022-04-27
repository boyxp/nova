package nova

import "os"
import "log"
import "runtime"
import "syscall"
import "strconv"
import "net/http"
import "io/ioutil"
import "github.com/fvbock/endless"
import "github.com/boyxp/nova/router"
import "github.com/boyxp/nova/request"
import "github.com/boyxp/nova/response"
import "github.com/boyxp/nova/register"
import "github.com/boyxp/nova/exception"
import _ "github.com/joho/godotenv/autoload"

func Listen(port string) *App {
	if port=="" {
		log.Printf("\033[1;31;40m%s\033[0m\n",".env配置文件不存在或port未设置")
		os.Exit(1)
	}
	return &App{port, &request.Form{}, &response.Json{}}
}

type App struct {
	Port string
	Request request.Interface
	Response response.Interface
}

func (A *App) Run() {
	server  := endless.NewServer("localhost:"+A.Port, A)
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		log.Println("pid:",pid)
		con := []byte(strconv.Itoa(pid))
		err := ioutil.WriteFile("pid", con, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	log.Println("Server stopped")

	os.Exit(0)
}

func (A *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer register.Clean()
	defer A.Catch()

	register.SetResponseWriter(w)
	register.SetRequest(r)

	match := router.Match(r.RequestURI)
	if match != true {
		exception.New("路由地址错误", 100)
	}

	params := A.Request.Parse(r)
	result := router.Invoke(r.RequestURI, params)
	if result != nil {
		A.Response.Render(result)
	}
}

func (A *App) SetRequest(req request.Interface) *App {
	A.Request = req
	return A
}

func (A *App) SetResponse(res response.Interface) *App {
	A.Response = res
	return A
}
/*
func (A *App) Before(route string, func(w http.ResponseWriter, r *http.Request)) *App {

}

func (A *App) After(route string, func(w http.ResponseWriter, r *http.Request)) *App {

}
*/
//异常捕获
func (A *App) Catch() {
        if err :=recover();err !=nil {
        		//断言逻辑异常直接抛出给用户
        		exception, ok := err.(*exception.Exception)
        		if ok {
        			A.Response.Error(exception.GetMessage(), exception.GetCode())

        			log.Println("逻辑异常代码：", exception.GetCode(), "逻辑异常内容：", exception.GetMessage())
        			return
        		}

        		//其他异常返回用户模糊提示
        		A.Response.Error("系统异常请联系管理员", -100)

        		//写入精确异常日志
        		log.Println("系统异常代码：-100","系统异常内容：", err)

                for i := 0; ; i++ {
                    _, file, line, ok := runtime.Caller(i)
                    if !ok {
                        break
                    }
                   	log.Println(i, file, line)
                }

                log.Println("\n")
        }
}
