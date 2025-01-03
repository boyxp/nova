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
import "github.com/joho/godotenv"

func init() {
	log.Println("nova init...")
	godotenv.Overload()
}

func Run() {
	Listen().Run()
}

func Listen(port ...string) *App {
	var Port string = "9800"
	if len(port)==0 {
		envPort := os.Getenv("port")
		if envPort=="" {
			log.Printf("\033[1;31;40m%s\033[0m\n",".env 配置文件不存在或port未设置,采用默认端口9800")
		} else {
			Port = envPort
		}
	} else {
		Port = port[0]
	}

	return &App{Port:Port, Request:&request.Form{}, Response:&response.Json{}}
}

func Register(controller interface{}) bool {
	return router.Register(controller)
}


type App struct {
	Port string
	Request request.Interface
	Response response.Interface
	middleware []func(next http.Handler) http.Handler
}

func (A *App) Use(next func(next http.Handler) http.Handler) *App {
	A.middleware = append(A.middleware, next)
	return A
}

func (A *App) Handle() http.Handler {
   var handler http.Handler = A

   for i:=len(A.middleware)-1;i>=0;i-- {
      handler = A.middleware[i](handler)
   }

   return handler
}

func (A *App) Run() {
	server  := endless.NewServer(":"+A.Port, A.Handle())
	server.BeforeBegin = func(add string) {
		pid := syscall.Getpid()
		con := []byte(strconv.Itoa(pid))
		err := ioutil.WriteFile("pid", con, 0644)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Process ID:", pid)
		log.Println("Listening and serving HTTP on:", A.Port)
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	log.Println("Server stopped")

	os.Exit(0)
}

func (A *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)

	defer register.Clean()
	defer A.Catch()

	register.SetResponseWriter(w)
	register.SetRequest(r)

	params := A.Request.Parse(r)
	result := router.Invoke(r.URL.Path, params)
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
        		log.Println("\033[31m系统异常代码：-100","系统异常内容：", err, "\033[0m")

                for i := 2; ; i++ {
                    _, file, line, ok := runtime.Caller(i)
                    if !ok {
                        break
                    }
                   	log.Println("\t", i-2, ")", file, line)
                }

                log.Println()
        }
}
